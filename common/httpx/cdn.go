package httpx

import (
	"fmt"
	"net"
	"github.com/miekg/dns"
)

var nameserver_all = []string{"1.1.1.1", // 美国 APNIC&CloudFlare公共DNS服务器
	"8.8.8.8",   // 美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器
	"77.88.8.8", // 俄国 搜索引擎 Yandex 服务器

	"61.47.33.9", // 新加坡 Pacific互联网泰国有限公司新加坡节点DNS服务器

	"202.14.67.4", // 香港 亚太环通(Pacnet)有限公司DNS服务器
	"61.10.1.130", // 香港 CableTVDNS服务器

	"118.118.118.118", // 上海市 电信DNS服务器(全国通用)
	"202.98.192.67",   // 贵州电信 DNS

	"58.132.8.1",      // 北京市 教育网DNS服务器
	"114.114.114.114", // 江苏省南京市 南京信风网络科技有限公司GreatbitDNS服务器
	"223.5.5.5",       // 浙江省杭州市 阿里巴巴anycast公共DNS
}

// CdnCheck verifies if the given ip is part of Cdn ranges
func (h *HTTPX) CdnCheck(ip string,host string) (bool, string, error) {
	if h.cdn == nil {
		return false, "", fmt.Errorf("cdn client not configured")
	}
	b, s, err := h.cdn.Check(net.ParseIP((ip)))
	if !b {
		return nslookup(host)
	}
	return b, s, err
	//return h.cdn.Check(net.ParseIP((ip)))
}

func nslookup(target string) (bool, string, error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)

	var res []string
	for _, dns_server := range nameserver_all {
		ns := dns_server + ":53"
		r, _, err := c.Exchange(&m, ns)
		if err != nil {
			//gologger.Error().Msgf("nameserver %s error: %v\n", ns, err)
			continue
		}
		//gologger.Info().Msgf("nameserver %s took %v\n", ns, t)
		if len(r.Answer) == 0 {
			continue
		}
		for _, ans := range r.Answer {
			record, isType := ans.(*dns.A)
			if isType {
				res = append(res, record.A.String())
			}
		}
	}

	res = RemoveDuplicatesAndEmpty(res)
	if len(res) > 1 {
		//justString := strings.Join(res, " ")
		return true, "CDN", nil
	}
	return false, "---", nil
}

//去除重复字符串和空格
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
