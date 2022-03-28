// test for issue: https://github.com/apache/dubbo-go-hessian2/issues/311

package issue311

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	hessian2 "github.com/apache/dubbo-go-hessian2"
	dubbo_hessian "github.com/apache/dubbo-go/protocol/dubbo/hessian2"
	"github.com/google/gopacket/pcapgo"
)

func loadPcap(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, _ := pcapgo.NewReader(f)
	packet, _, err := r.ReadPacketData()
	if err != nil || packet == nil {
		return nil, err
	}
	return packet, err
}

// dubbo 协议参考： https://zhuanlan.zhihu.com/p/98562180
func TestParseDubboHessian(t *testing.T) {
	raw, _ := loadPcap("dubbo_req.pcap")

	codecR := hessian2.NewHessianCodec(bufio.NewReader(bytes.NewReader(raw[108:])))
	header := &hessian2.DubboHeader{}

	// 从二进制数据解析出的信息：
	// target: cn.com.xib.ifsp.loan.provider.xib.IXibQueryApplyProgressListProvider
	// serviceVersion: 1.0.0
	// method: $invoke
	// dubboVersion: 2.8.4a
	// argsTypes: 解析有问题，解析出应该为string，但是int值-1
	err := codecR.ReadHeader(header)
	if err != nil {
		t.Errorf("read header err [%v]\n", err)
	}

	c := make([]interface{}, 7)
	err = codecR.ReadBody(c)
	if err != nil {
		t.Errorf("read body err [%v]\n", err)
	}
}

func TestParseDubboGoHessian(t *testing.T) {
	raw, _ := loadPcap("dubbo_req.pcap")

	codecR := dubbo_hessian.NewHessianCodec(bufio.NewReader(bytes.NewReader(raw[108:])))
	header := &dubbo_hessian.DubboHeader{}
	err := codecR.ReadHeader(header)
	if err != nil {
		t.Errorf("read header err [%v]\n", err)
	}

	c := make([]interface{}, 7)
	err = codecR.ReadBody(c)
	if err != nil {
		t.Errorf("read body err [%v]\n", err)
	}
}
