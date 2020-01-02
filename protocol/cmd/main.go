//用于更加注册文件自动注册pb
package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	gAllItems    []*pbItem
	gAppendItems []*pbItem
	gItemMap     = make(map[string]uint16)
	gMaxMsgId    uint16
)

type pbItem struct {
	id    uint16
	msg   string
	comma string
}

func main() {
	regText := "proto/register.txt"
	pbDir := "pb"
	if len(os.Args) > 1 {
		sep := string(os.PathSeparator)
		regText = os.Args[1] + sep + regText
		pbDir = os.Args[1] + sep + pbDir
	}
	loadRegisterFile(regText)
	loadProtoDir(pbDir)
	genFile()
	writerNewProto(regText)
}

// 读取已注册文件
func loadRegisterFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open input file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	var idMap = make(map[uint16]string)
	for {
		s, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(s)
		commaStr := ""
		idx := strings.Index(line, "//")
		if idx >= 0 {
			commaStr = line[idx:]
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		its := strings.Split(line, "=")
		if len(its) != 2 {
			fmt.Printf("Invalid line %v\n", its)
			break
		}
		sMsgId := strings.ToLower(strings.TrimSpace(its[1]))
		base := 10
		if strings.HasPrefix(sMsgId, "0x") {
			sMsgId = sMsgId[2:]
			base = 16
		}

		msgId, err := strconv.ParseUint(sMsgId, base, 16)
		if err != nil {
			fmt.Printf("ParseUint err at line %s\n %s", line, string(s))
		}
		msgStr := strings.TrimSpace(its[0])
		if _, ok := idMap[uint16(msgId)]; ok {
			fmt.Printf("register twice for msgId: %d\n", msgId)
			os.Exit(1)
		}
		id := uint16(msgId)
		idMap[id] = msgStr
		gAllItems = append(gAllItems, &pbItem{
			id:    id,
			msg:   msgStr,
			comma: commaStr,
		})
		gItemMap[msgStr] = id
		if id > gMaxMsgId {
			gMaxMsgId = id
		}
	}
}

//读取协议目录
func loadProtoDir(dirPath string) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	sep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			subDir := dirPath + sep + fi.Name()
			loadProtoDir(subDir)
		} else {
			ext := strings.ToLower(filepath.Ext(fi.Name()))
			if ext == ".go" {
				loadProtoFile(dirPath + sep + fi.Name())
			}
		}
	}
}

// 读取协议
func loadProtoFile(fileName string) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		ret, ok := n.(*ast.GenDecl)
		if ok {
			for _, spec := range ret.Specs {
				sp, ok := spec.(*ast.TypeSpec)
				if ok {
					msgName := sp.Name.String()
					if strings.HasPrefix(msgName, "Message") {
						_, ok = gItemMap[msgName]
						if !ok {
							gMaxMsgId++
							it := &pbItem{
								id:    gMaxMsgId,
								msg:   msgName,
								comma: ret.Doc.Text(),
							}
							gAllItems = append(gAllItems, it)
							gAppendItems = append(gAppendItems, it)
							gItemMap[it.msg] = it.id
						}
					}
					return true
				}
			}
		}
		return true
	})
}

// 生成go文件
func genFile() {
	f, err := os.Create("protoFactory.go")
	if err != nil {
		fmt.Printf("genFile error:%v", err)
		os.Exit(1)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	bw.WriteString("// Code generated by go proto generate tool. DO NOT EDIT.\n")
	bw.WriteString("package protocol\n")
	bw.WriteString("import \"gogame/protocol/pb\"\n")
	bw.WriteString("const(\n")
	for _, it := range gAllItems {
		bw.WriteString(fmt.Sprintf("MsgID_%s = %d %s\n", it.msg, it.id, it.comma))
	}
	bw.WriteString(")\n")

	bw.WriteString("func init() {\n")
	bw.WriteString("processor = NewProcessor()\n")
	for _, it := range gAllItems {
		bw.WriteString(fmt.Sprintf("\tprocessor.Register(MsgID_%s, (*pb.%s)(nil))\n", it.msg, it.msg))
	}
	bw.WriteString("}\n")
	bw.Flush()
	fmt.Println("genFile Ok ")
}

//回写新增的协议
func writerNewProto(fileName string) {
	if len(gAppendItems) == 0 {
		return
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("open input file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()
	bw.WriteString("\n")
	for _, it := range gAppendItems {
		bw.WriteString(fmt.Sprintf("%s = 0x%X // %s\n", it.msg, it.id, it.comma))
	}
	fmt.Printf("writer %v Ok\n", fileName)
}
