package httpx

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
)

var wafProbe = "/?id=1 AND 1=1 UNION ALL SELECT 1,NULL,'<script>alert(\"XSS\")</script>',table_name FROM information_schema.tables WHERE 2>1--/**/; EXEC xp_cmdshell('cat ../../../etc/passwd')"

func (h *HTTPX) WafCheck(host string) (bool, string, error) {
	return WafGet(host)
}

func WafGet(host string) (bool, string, error) {
	respWafProbe, err := http.Get("https://" + host + wafProbe)
	respWafProbe.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	respWafProbe.Header.Set("Referer", "host")
	if err != nil {
		panic(err)
	}
	defer respWafProbe.Body.Close()
	respWafProbeBody, err := ioutil.ReadAll(respWafProbe.Body)
	if err != nil {
		fmt.Printf("请求报错")
	}
	respWafProbeBodyStr := string(respWafProbeBody)
	respWafProbeHeader, _ := json.Marshal(respWafProbe.Header)
	respWafProbeHeaderStr := string(respWafProbeHeader)
	textProbe := respWafProbeHeaderStr + respWafProbeBodyStr
	//log.Println(respWafProbeBodyStr)
	//log.Println(textProbe)
	wafdata, err := simplejson.NewJson([]byte(wafdatajson))
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, v := range wafdata.Get("wafs").MustMap(){
		i,_ := v.(map[string]interface{})
		if i["regex"].(string) == "" {
			continue
		}
		re, err := regexp.Compile(i["regex"].(string))
		if err != nil {
			log.Fatal(err)
		}
		found := re.MatchString(textProbe)
		//matched, err := regexp.Match(i["regex"].(string), []byte(respWafProbeBodyStr))
		if found {
			//log.Println(i["regex"].(string),err)
			return true, i["name"].(string), nil
		}
	}
	return false, "---", nil
}

var wafdatajson string = `
{
    "__copyright__": "Copyright (c) 2019-2021 Miroslav Stampar (@stamparm), MIT. See the file 'LICENSE' for copying permission",
    "__notice__": "The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software",

    "payloads": [
        "HTML::<img>",
        "SQLi::1 AND 1",
        "SQLi::1/**/AND/**/1",
        "SQLi::1/*0AND*/1",
        "SQLi::1 AND 1=1",
        "SQLi::1 AND 1 LIKE 1",
        "SQLi::1 AND 1 BETWEEN 0 AND 1",
        "SQLi::1 AND 2>(SELECT 1)-- -",
        "SQLi::' OR SLEEP(5) OR '",
        "SQLi::admin'-- -",
        "SQLi::information_schema",
        "SQLi::;DROP TABLE mysql.users",
        "SQLi::';DROP DATABASE mysql#",
        "SQLi::1/**/UNION/**/SELECT/**/1/**/FROM/**/information_schema.*",
        "SQLi::SELECT id FROM users WHERE id>2",
        "SQLi::1 UNION SELECT information_schema.*",
        "SQLi::1;EXEC xp_cmdshell('type autoexec.bat');",
        "SQLi::1;INSERT INTO USERS values('admin', 'foobar')",
        "XSS::<img src=x onerror=alert('XSS')>",
        "XSS::<img onfoo=f()>",
        "XSS::<script>",
        "XSS::<script>alert('XSS')</script>",
        "XSS::\\\";alert('XSS');//",
        "XSS::1' onerror=alert(String.fromCharCode(88,83,83))>",
        "XSS::<![CDATA[<script>var n=0;while(true){n++;}</script>]]>",
        "XSS::<meta http-equiv=\"refresh\" content=\"0;url=data:text/html;base64,PHNjcmlwdD5hbGVydCgnWFNTJyk8L3NjcmlwdD4K\">",
        "XSS::javascript:alert(/XSS/)",
        "XSS::<marquee onstart=alert(1)>",
        "XPATHi::' and count(/*)=1 and '1'='1",
        "XPATHi::count(/child::node())",
        "XPATHi::' and count(/comment())=1 and '1'='1",
        "XPATHi::' or '1'='1",
        "XXE::<!ENTITY xxe SYSTEM \"file:///etc/passwd\" >]><foo>&xxe;</foo>",
        "LDAPi::admin*)((|userpassword=*)",
        "LDAPi::user=*)(uid=*))(|(uid=*",
        "LDAPi::*(|(objectclass=*))",
        "NOSQLi::true, $where: '1 == 1'",
        "NOSQLi::{ $ne: 1 }",
        "NOSQLi::' } ], $comment:'success'",
        "PHPi::<?php include_once(\"/etc/passwd\"); ?>",
        "ACE::netstat -antup | grep :443; ping 127.0.0.1; curl http://www.google.com",
        "PT:://///.htaccess",
        "PT::/etc/passwd",
        "PT::../../boot.ini",
        "PT::C:/inetpub/wwwroot/global.asa"
    ],
    "wafs": {
        "360": {
            "company": "360",
            "name": "360",
            "regex": "<title>493</title>|/wzws-waf-cgi/",
            "signatures": [
                "9778:RVZXum61OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTC",
                "9ccc:RVZXum61OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "aesecure": {
            "company": "aeSecure",
            "name": "aeSecure",
            "regex": "aesecure_denied\\.png|aesecure-code: \\d+",
            "signatures": [
                "8a4b:RVdXu260OEhCWapBYKcPk4JzWOtohM4JiUcMrmRXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJOdLsXo2tKaK99n+i7c4RmkgI2FZnxtDtBeq+c36A4chW1XaTD"
            ]
        },
        "airlock": {
            "company": "Phion/Ergon",
            "name": "Airlock",
            "regex": "The server detected a syntax error in your request",
            "signatures": [
                "3e2c:RVZXu261OEhCWapBYKcPk4JzWOtohM4IiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXomtKaK59n+i6c4RmkwI2FZjxtDtAeq6c36A5chW1XaTD"
            ]
        },
        "alertlogic": {
            "company": "Alert Logic",
            "name": "Alert Logic",
            "regex": "(?s)timed_redirect\\(seconds, url\\).+?<p class=\"lid\">Reference ID:",
            "signatures": []
        },
        "aliyundun": {
            "company": "Alibaba Cloud Computing",
            "name": "AliYunDun",
            "regex": "Sorry, your request has been blocked as it may cause potential threats to the server's security|//errors\\.aliyun\\.com/",
            "signatures": [
                "e082:RVZXum61OElCWapAYKYPkoJzWOpohM4JiUYMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC"
            ]
        },
        "anquanbao": {
            "company": "Anquanbao",
            "name": "Anquanbao",
            "regex": "/aqb_cc/error/",
            "signatures": [
                "c790:RVZXum61OElCWapAYKYPk4JzWOpohM4JiUYMr2RXg1uQJbX3uhdOn9hsOj+hXrAB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "d3d3:RVZXum61OElCWapAYKYPk4JzWOpohM4JiUYMr2RXg1uQJbX3uhdOn9hsOj+hXrAB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC"
            ]
        },
        "approach": {
            "company": "Approach",
            "name": "Approach",
            "regex": "Approach.+?Web Application (Firewall|Filtering)",
            "signatures": [
                "fef0:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A5chW1XKTD"
            ]
        },
        "armor": {
            "company": "Armor Defense",
            "name": "Armor Protection",
            "regex": "This request has been blocked by website protection from Armor",
            "signatures": [
                "03ec:RVZXum60OEhCWapBYKYPk4JzWOtohM4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c36A4chS1XaTC",
                "1160:RVZXum60OEhCWapBYKYPk4JyWOtohM4IiUcMr2RWg1qQJbX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ],
            "note": "Uses SecureSphere (Imperva) (Reference: https://www.imperva.com/resources/case_studies/CS_Armor.pdf)"
        },
        "asm": {
            "company": "F5 Networks",
            "name": "Application Security Manager",
            "regex": "The requested URL was rejected\\. Please consult with your administrator|security\\.f5aas\\.com",
            "signatures": [
                "2f81:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI3FZjxtDtAeq+c36A4chS1XaTC",
                "4fd0:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq6c3qA4chS1XaTC",
                "5904:RVZXum60OEhCWapBYKcPk4JzWOpohc4IiUcMr2RWg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtTtAeq+c3qA4chS1XaTC",
                "8bcf:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtTtAeq6c36A5chS1XaTC",
                "540f:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtTtAeq+c36A5chS1XaTC",
                "c7ba:RVZXum60OEhCWKpAYKYPkoJzWOpohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtLaK99n+i7c4VmkwI3FZjxtDtAeq6c3qA4chS1XaTC",
                "fb21:RVZXum60OEhCWapBYKcPk4JzWOpohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI3FZjxtDtAeq+c36A5chW1XaTC",
                "b6ff:RVZXum61OEhCWapBYKcPkoJzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq+c36A4chW1XaTC",
                "3b1e:RVZXum60OEhCWapBYKcPk4JyWOpohM4IiUcMr2RWg1qQJLX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tKaK99nui7c4RmkgI2FZjxtDtAeq6c3qA5chS1XKTC",
                "620c:RVZXum60OEhCWapBYKcPkoJzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC",
                "b9a0:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq+c3qA4chW1XaTC",
                "ccb6:RVdXum61OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtTtAeq+c36A5chW1XaTC",
                "9138:RVZXum60OEhCWapBYKcPk4JzWOpohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq6c3qA4chS1XaTC",
                "54cc:RVZXum61OEhCWapBYKcPkoJzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq6c3qA4chS1XaTC",
                "4c83:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTC",
                "8453:RVZXum60OEhCWapBYKcPk4JzWOtohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI3FZjxtDtAeq+c36A4chS1XaTC"
            ]
        },
        "astra": {
            "company": "Czar Securities",
            "name": "Astra",
            "regex": "(?s)unfortunately our website protection system.+?//www\\.getastra\\.com",
            "signatures": []
        },
        "aws": {
            "company": "Amazon",
            "name": "AWS WAF",
            "regex": "(?i)HTTP/1.+\\b403\\b.+\\s+Server: aws|(?s)Request blocked.+?Generated by cloudfront",
            "signatures": [
                "2998:RVZXu261OEhCWapBYKcPk4JzWOpohM4IiUcMr2RWg1uQJbX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "fffa:RVZXum60OEhCWapAYKYPk4JyWOpohc4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "9de0:RVZXu261OEhCWapBYKcPk4JzWOpohM4IiUcMr2RWg1uQJbX3uhZOnthtOj+hXrAA16BcPhJOdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "34a8:RVZXu261OEhCWapBYKcPk4JzWOpohM4IiUcMr2RWg1uQJbX3uhdOn9htOj+hXrAB16BcPxJOdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "1104:RVZXum61OEhCWapBYKcPk4JzWOpohM4IiUcMr2RXg1uQJbX3uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "ea40:RVZXu261OEhCWapBYKcPk4JzWOtohM4IiUcMr2RWg1uQJbX3uhdOn9htOj+hXrAB16BcPxJOdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "barracuda": {
            "company": "Barracuda Networks",
            "name": "Barracuda",
            "regex": "\\bbarracuda_|barra_counter_session=|when this page occurred and the event ID found at the bottom of the page|<!--(0123456789){15}",
            "signatures": [
                "2676:RVdXum61OElCWapAYKYPk4JzWOtohM4JiUcMr2RWg1qQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXo2tKaK99n+i6c4VmkwI3FZjxtDtAeq6c36A4chS1XaTC",
                "db27:RVdXum61OElCWapAYKYPk4JzWOtohM4JiUcMr2RWg1qQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XaTC",
                "7e6b:RVdXum61OElCWapBYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A4chS1XaTC"
            ]
        },
        "bekchy": {
            "company": "Faydata Information Technologies Inc.",
            "name": "Bekchy",
            "regex": "<title>Bekchy - Access Denided</title>|<a class=\"btn\" href=\"https://bekchy.com/report\">",
            "signatures": [
                "e1c5:RVZXum60OEhCWKpAYKYPk4JzWOtohc4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "bitninja": {
            "company": "BitNinja",
            "name": "BitNinja",
            "regex": "alt=\"BitNinja|Security check by BitNinja|your IP will be removed from BitNinja|<title>Visitor anti-robot validation</title>",
            "signatures": []
        },
        "bluedon": {
            "company": "Bluedon",
            "name": "Bluedon",
            "regex": "Bluedon Web Application Firewall|Server: BDWAF",
            "signatures": []
        },
        "bulletproof": {
            "company": "AITpro Website Security",
            "name": "BulletProof Security Pro",
            "regex": "(?s)bpsMessage.+?403 Forbidden Error Page.+?If you arrived here due to a search or clicking on a link",
            "signatures": []
        },
        "cdnns": {
            "company": "CdnNs/WdidcNet",
            "name": "CdnNsWAF",
            "regex": "by CdnNsWAF Application Gateway",
            "signatures": [
                "5c5d:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXo2tLaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC"
            ]
        },
        "cerber": {
            "company": "Cerber Tech",
            "name": "WP Cerber Security",
            "regex": "We're sorry, you are not allowed to proceed|Your request looks suspicious or similar to automated requests from spam posting software",
            "signatures": [
                "d8c2:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "checkpoint": {
            "company": "Check Point",
            "name": "Next Generation Firewall",
            "regex": "",
            "signatures": [
                "b771:RVZXum61OEhCWapAYKYPkoJzWOpohc4JiUYMr2RWg1uQJbX2uhdOnthsOj+hX7AB16BcPhJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "3b40:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUYMrmRWg1qQJLX2uhdOnthsOj+hX7AB16BcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XKTC",
                "a332:RVZXum61OEhCWapAYKYPkoJzWOpohc4JiUYMr2RWg1uQJbX2uhdOnthsOj+hX7AB16BcPhJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "a89b:RVZXum61OEhCWapAYKYPkoJzWOpohc4JiUYMr2RWg1uQJbX2uhdOnthsOj+hX7AB16BcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC"
            ]
        },
        "chuangyu": {
            "company": "Yunaq",
            "name": "Chuang Yu Shield",
            "regex": " \\d+\\.\\d+\\.\\d+\\.\\d+/[0-9a-f]{7} \\[\\d+\\] ",
            "signatures": [
                "eda6:RVZXum61OElCWapAYKcPkoJzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTC",
                "5bae:RVZXum61OElCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC"
            ]
        },
        "cloudbric": {
            "company": "Cloudbric",
            "name": "Cloudbric",
            "regex": "Your request was blocked by Cloudbric",
            "signatures": [
                "514d:RVZXum60OEhCWapBYKcPk4JzWOtohM4JiUcMrmRXg1qQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "cloudflare": {
            "company": "CloudFlare",
            "name": "CloudFlare",
            "regex": "Attention Required! \\| Cloudflare|CLOUDFLARE_ERROR_",
            "signatures": [
                "956d:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "6b42:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "2295:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "0d86:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhdOnthsOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "4849:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMrmRWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "535c:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUYMr2RWg1uQJbX2uhdOnthtOj+hXrAB16FcPxJOdLoXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "675a:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMrmRWg1uQJbX2uhdOnthsOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "4a45:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUcMrmRWg1uQJLX2uhdOnthsOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC",
                "1f29:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUcMrmRWg1uQJLX2uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "6002:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUcMrmRWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "78df:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMrmRWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTD",
                "cf65:RVZXum60OEhCWapBYKcPkoJzWOtohM4IiUcMrmRWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4VmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "85c6:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC",
                "9a2d:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMrmRWg1uQJLX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "0576:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMrmRXg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "f3bb:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUYMr2RXg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "471d:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "8936:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUcMrmRWg1uQJLX2uhdOnthsOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC",
                "0ade:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "22d1:RVZXum60OEhCWapBYKcPkoJzWOpohM4IiUcMr2RWg1uQJbX2uhdOnthtOj+hXrAA16FcPxJOdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "e9bd:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthsOj+hXrAB16FcPxJPdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "comodo": {
            "company": "Comodo",
            "name": "Comodo",
            "regex": "Server: Protected by COMODO WAF",
            "signatures": [
                "ade8:RVZXum60OEhCWapAYKYPkoJzWOpohc4IiUYMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTD",
                "f063:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTD",
                "985c:RVZXum60OEhCWapAYKYPkoJzWOpohc4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c3qA5chW1XaTD",
                "f063:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTD",
                "1971:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD"
            ]
        },
        "crawlprotect": {
            "company": "Jean-Denis Brun",
            "name": "CrawlProtect",
            "regex": "<title>CrawlProtect|This site is protected by CrawlProtectc|Set-Cookie: crawlprotecttag",
            "signatures": [
                "1eca:RVZXum60OEhCWKpBYKYPkoJzWOpohM4IiUYMrmRXg1uQJLX2uhZOnthtOj+hXrAA16FcPhJPdLoXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XKTC"
            ]
        },
        "distil": {
            "company": "Distil Networks",
            "name": "Distil",
            "regex": "distilCaptchaForm|distilCallbackGuard|cdn\\.distilnetworks\\.com/images/anomaly-detected\\.png",
            "signatures": []
        },
        "dotdefender": {
            "company": "Applicure Technologies",
            "name": "dotDefender",
            "regex": "dotDefender Blocked Your Request|Applicure is the leading provider of web application security|Please contact the site administrator, and provide the following Reference ID",
            "signatures": [
                "7cce:RVZXum60OEhCWapAYKYPkoJzWOpohM4IiUYMrmRWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "dddb:RVdXum61OElCWapAYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "0718:RVZXum61OElCWapAYKYPk4JzWOtohM4IiUYMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "9bf2:RVdXum61OElCWapAYKYPk4JzWOtohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC"
            ]
        },
        "duedge": {
            "company": "Baidu",
            "name": "DuEdge",
            "regex": "(?s)<h1>403<small>.+DuEdge Event ID: [0-9a-f]{16}.+IP: [0-9.]+",
            "signatures": []
        },
        "expressionengine": {
            "company": "EllisLab",
            "name": "ExpressionEngine",
            "regex": "(?s)\\bexp_last_.+?(Invalid GET Data|Invalid URI)",
            "signatures": [
                "88ec:RVZXum60OEhCWKpAYKYPkoJyWOpohM4JiUcMrmRWg1qQJbX3uhZOnthsOj6hX7AA16FcPxJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chS1XKTC"
            ]
        },
        "fortiweb": {
            "company": "Fortinet",
            "name": "FortiWeb",
            "regex": "Server Unavailable!",
            "signatures": [
                "9d05:RVZXu261OElCWapBYKcPk4JzWOtohM4IiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI3FZjxtDtAeq+c36A5chW1XaTD"
            ]
        },
        "godaddy": {
            "company": "GoDaddy",
            "name": "GoDaddy Website Security",
            "regex": "GoDaddy Security - Access Denied|Access Denied - GoDaddy Website Firewall",
            "signatures": [
                "6cff:RVdXum60OEhCWapAYKYPk4JzWOtohM4IiUYMr2RWg1uQJbX3uhdOn9htOj+hXrAA16FcPxJOdLoXomtKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "greywizard": {
            "company": "Grey Wizard",
            "name": "Greywizard",
            "regex": "(?i)server: greywizard|detected attempted attack or non standard traffic from your IP address|<title>Grey Wizard</title>",
            "signatures": [
                "c669:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOnthsOj+hX7AB16FcPhJPdLsXomtKaK59nui7c4RmkwI2FZjxtDtAeq+c3qA5chW1XaTC"
            ]
        },
        "gtmc": {
            "company": "GTMC",
            "name": "GTMC WAF",
            "regex": "GTMC WAF1 Protection:|Please consult with administrator or waf@nm.gtmc.com.tw",
            "signatures": []
        },
        "imunify360": {
            "company": "CloudLinux",
            "name": "Imunify360",
            "regex": "Server: imunify360-webshield|protected by Imunify360|Powered by Imunify360|imunify360 preloader",
            "signatures": []
        },
        "incapsula": {
            "company": "Incapsula/Imperva",
            "name": "Incapsula",
            "regex": "Incapsula incident ID",
            "signatures": [
                "2770:RVZXum60OEhCWKpAYKYPkoJzWOpohc4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC",
                "3193:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZnxtDtAeq6c3qA4chS1XKTC",
                "cdd1:RVZXum60OEhCWapAYKcPk4JzWOpohM4IiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtLaK99n+i7c4RmkgI2FZnxtTtBeq+c36A5chW1XaTC"
            ]
        },
        "isaserver": {
            "company": "Microsoft",
            "name": "ISA Server",
            "regex": "The (ISA Server|server) denied the specified Uniform Resource Locator \\(URL\\)",
            "signatures": []
        },
        "ithemes": {
            "company": "iThemes",
            "name": "iThemes Security",
            "regex": "",
            "signatures": [
                "c70f:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "71ee:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1qQJLX2uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ],
            "note": "Formerly Better WP Security"
        },
        "janusec": {
            "company": "Janusec",
            "name": "Janusec Application Gateway",
            "regex": "Reason:.+by Janusec Application Gateway",
            "signatures": [
                "5c5d:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXo2tLaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC"
            ]
        },
        "jiasule": {
            "company": "Jiasule",
            "name": "Jiasule",
            "regex": "Server: jiasule-WAF|notice-jiasule|static\\.jiasule\\.com/static/js/http_error\\.js",
            "signatures": [
                "7520:RVZXum61OElCWapAYKYPk4JzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtBeq+c36A5chW1XaTD",
                "001e:RVZXum61OElCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI3FZjxtTtAeq+c36A5chW1XaTC",
                "665d:RVZXum61OElCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA5chS1XaTC",
                "4fed:RVZXum61OElCWapAYKYPkoJzWOpohM4IiUYMr2RXg1uQJbX2uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC"
            ]
        },
        "knownsec": {
            "company": "Knownsec",
            "name": "KS-WAF",
            "regex": "url\\('/ks-waf-error\\.png'\\)",
            "signatures": []
        },
        "kona": {
            "company": "Akamai Technologies",
            "name": "Kona Site Defender",
            "regex": "(?s)Server: AkamaiGHost.+?You don't have permission to access|\\b18\\.[0-9a-f]{8}.1[0-9]{9}\\.[0-9a-f]{7}\\b",
            "signatures": [
                "b996:RVZXum60OEhCWapAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "1893:RVZXum60OEhCWapAYKYPk4JzWOtohM4JiUcMr2RXg1uQJLX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkwI2FZjxtDtAeq+c3qA4chS1XKTC",
                "165b:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "12b3:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "3426:RVZXum60OEhCWapAYKYPk4JzWOtohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "e197:RVZXum60OEhCWKpAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "eb57:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hX7AB16FcPxJPdLsXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c36A4chS1XaTC",
                "94ed:RVZXum60OEhCWapAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "5ca8:RVZXum60OEhCWKpAYKYPkoJzWOtohM4IiUYMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXomtKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "cc5b:RVZXum60OEhCWKpAYKYPkoJzWOtohM4IiUYMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "e7d9:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "bd78:RVZXum60OEhCWKpAYKYPk4JzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "6cbc:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD",
                "a40d:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "1f03:RVZXum60OEhCWapBYKYPk4JzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD",
                "e120:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "7ae5:RVZXum60OEhCWKpAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "6bf2:RVZXum60OEhCWapAYKYPkoJzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "1db3:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "fcbb:RVZXum60OEhCWapAYKYPkoJzWOtohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "d1b6:RVZXum60OEhCWKpAYKYPkoJzWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "8b30:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "8db8:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "8900:RVZXum60OEhCWapAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "677e:RVZXum60OEhCWapAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "a13a:RVZXum60OEhCWKpAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "579e:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "82b4:RVZXum60OEhCWapAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD",
                "22e4:RVZXum60OEhCWapAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "bd0e:RVZXum60OEhCWapAYKYPk4JzWOtohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "8976:RVZXum60OEhCWKpAYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "e34c:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1qQJLX2uhdOn9htOj+hX7AB16FcPxJPdLsXomtKaK59nui6c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC"
            ]
        },
        "kuipernet": {
            "company": "ASTSoft",
            "name": "Kuipernet",
            "regex": "(?s)Content-Length: 118214.+W5M0MpCehiHzreSzNTczkc9d",
            "signatures": []
        },
        "malcare": {
            "company": "Inactiv",
            "name": "MalCare",
            "regex": "Blocked because of Malicious Activities|Firewall(<[^>]+>)*powered by(<[^>]+>)*MalCare",
            "signatures": [
                "def2:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "modsecurity": {
            "company": "Trustwave",
            "name": "ModSecurity",
            "regex": "(?i)Server:.+mod_security|This error was generated by Mod_Security|/modsecurity\\-errorpage/|One or more things in your request were suspicious|rules of the mod_security module|mod_security rules triggered|Protected by Mod Security|HTTP Error 40\\d\\.0 - ModSecurity Action|40\\d ModSecurity Action|ModSecurity IIS \\(\\d+bits\\)</td>",
            "signatures": [
                "46d5:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "1ece:RVZXum61OEhCWapBYKcPk4JzWOpohc4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "69c6:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthsOj+hX7AB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "28eb:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthtOj+hXrAB16FcPhJOdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XaTC",
                "3918:RVZXum60OEhCWapAYKYPk4JyWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK99n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "511d:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "f694:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhZOnthtOj+hX7AB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "51ca:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJOdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "e18b:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhZOnthtOj+hX7AB16FcPhJOdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "6e99:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "dd72:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "f53e:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "e15c:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhZOnthtOj+hX7AB16FcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "ded8:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhZOnthtOj+hXrAB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "6e99:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "7986:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPhJOdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "02b2:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "4602:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJOdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "b1a2:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "5e9a:RVZXum60OEhCWapAYKYPk4JyWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hXrAB16FcPhJPdLsXomtKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "35c4:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chS1XKTC",
                "c697:RVZXum60OEhCWapAYKYPk4JyWOpohM4JiUcMr2RXg1uQJbX3uhZOnthtOj+hX7AB16FcPhJPdLsXomtKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "85e3:RVZXum60OElCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hX7AB16FcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "7d7f:RVZXum60OEhCWapAYKYPk4JyWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD",
                "064b:RVZXum60OEhCWapAYKYPk4JyWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hXrAB16FcPhJOdLsXomtKaK99n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "5659:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUYMr2RXg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "94b1:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "7951:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUcMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD",
                "b83a:RVZXum60OEhCWKpAYKYPkoJyWOpohM4JiUYMrmRWg1qQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTD",
                "4191:RVZXum60OEhCWapAYKYPkoJyWOpohM4JiUYMr2RXg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLoXomtKaK59n+i7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTD"
            ]
        },
        "naxsi": {
            "company": "NBS System",
            "name": "NAXSI",
            "regex": "(?i)Blocked By NAXSI|Naxsi Blocked Information|naxsi/waf",
            "signatures": [
                "19ee:RVdXum61OElCWKpAYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI3FZnxtDtBeq+c36A4chW1XaTC"
            ]
        },
        "netscaler": {
            "company": "Citrix",
            "name": "NetScaler AppFirewall",
            "regex": "<title>Application Firewall Block Page</title>|Violation Category: APPFW_|AppFW Session ID|Access has been blocked - if you feel this is in error, please contact the site administrators quoting the following",
            "signatures": [
                "9c6c:RVdXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMrmRWg1qQJbX3uhdOn9hsOj6hXrAA16BcPhJOdLsXo2tKaK99n+i6c4RmkgI2FZnxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "newdefend": {
            "company": "Newdefend",
            "name": "Newdefend",
            "regex": "Server: NewDefend|/nd_block/",
            "signatures": [
                "1ba1:RVZXu261OElCWapBYKYPk4JzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthsOj+hX7AB16FcPxJPdLoXo2tKaK99n+i7c4RmkwI3FZjxtDtAeq+c36A4chW1XaTD"
            ]
        },
        "nexusguard": {
            "company": "Nexusguard Limited",
            "name": "Nexusguard",
            "regex": "speresources\\.nexusguard\\.com/wafpage/[^>]*#\\d{3};|<p>Powered by Nexusguard</p>",
            "signatures": [
                "869d:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hX7AB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC"
            ]
        },
        "ninjafirewall": {
            "company": "NinTechNet",
            "name": "NinjaFirewall",
            "regex": "<title>NinjaFirewall: 403 Forbidden|For security reasons?, it was blocked and logged",
            "signatures": [
                "2c12:RVZXum60OEhCWapBYKYPkoJzWOtohM4JiUcMr2RXg1uQJLX3uhdOn9hsOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtBeq+c3qA4chW1XaTC"
            ]
        },
        "onmessageshield": {
            "company": "Blackbaud",
            "name": "onMessage Shield",
            "regex": "This site is protected by an enhanced security system to ensure a safe browsing experience|onMessage SHIELD",
            "signatures": [
                "125a:RVdXum61OElCWKpAYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI3FZnxtDtBeq+c36A5chW1XaTC"
            ]
        },
        "openrasp": {
            "company": "Blackbaud",
            "name": "OpenRASP",
            "regex": "400 - Request blocked by OpenRASP|https://rasp.baidu.com/blocked2?/",
            "signatures": []
        },
        "paloalto": {
            "company": "Palo Alto Networks",
            "name": "Palo Alto",
            "regex": "has been blocked in accordance with company policy|Palo Alto Next Generation Security Platform",
            "signatures": [
                "862a:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhZOnthsOj+hXrAA16BcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC",
                "5fe6:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMrmRWg1uQJLX2uhZOnthsOj+hXrAA16BcPhJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC",
                "cffd:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhZOnthsOj+hXrAA16BcPhJPdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC",
                "1427:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthtOj+hXrAA16FcPhJPdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "fa37:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhZOnthsOj6hXrAA16BcPhJOdLoXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "9135:RVZXum60OEhCWapAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX3uhZOnthsOj+hXrAA16BcPhJOdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC",
                "953a:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj+hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC"
            ]
        },
        "perimeterx": {
            "company": "PerimeterX",
            "name": "PerimeterX",
            "regex": "https://www.perimeterx.com/whywasiblocked",
            "signatures": []
        },
        "profense": {
            "company": "ArmorLogic",
            "name": "Profense",
            "regex": "Server: Profense",
            "signatures": [
                "eaee:RVZXum60OEhCWapAYKYPkoJyWOtohM4JiUcMr2RWg1uQJbX3uhdOnthsOj+hXrAB16FcPxJOdLsXo2tLaK99n+i6c4VmkwI3FZjxtDtAeq6c3qA4chS1XaTC"
            ]
        },
        "radware": {
            "company": "Radware",
            "name": "AppWall",
            "regex": "Unauthorized Request Blocked|You are seeing this page because we have detected unauthorized activity|mailto:CloudWebSec@radware\\.com",
            "signatures": [
                "e68e:RVdXu261OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXo2tKaK99n+i7c4VmkwI3FZnxtDtAeq+c36A5chW1XaTD",
                "48fa:RVdXu260OEhCWapBYKcPkoJzWOpohM4JiUYMrmRXg1uQJbX3uhdOn9hsOj+hX7AA16BcPxJOdLsXomtKaK59n+i6c4RmkgI2FZnxtDtAeq6c3qA5chW1XaTD",
                "8fc4:RVdXu261OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI3FZnxtDtAeq+c36A5chW1XaTD"
            ]
        },
        "reblaze": {
            "company": "Reblaze",
            "name": "Reblaze",
            "regex": "For further information, do not hesitate to contact us",
            "signatures": [
                "86fb:RVZXum61OElCWKpAYKcPkoJzWOtohM4JiUcMr2RXg1uQJbX3uhdOnthsOj6hXrAB16BcPhJPdLoXo2tLaK99n+i7c4RmkgI2FZjxtDtBeq+c36A5chW1XaTD"
            ]
        },
        "requestvalidationmode": {
            "company": "Microsoft",
            "name": "ASP.NET RequestValidationMode",
            "regex": "HttpRequestValidationException|Request Validation has detected a potentially dangerous client input value|ASP\\.NET has detected data in the request that is potentially dangerous",
            "signatures": [
                "7ecd:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hXrAA16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC",
                "919b:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hXrAA16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTD",
                "14fa:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hXrAA16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XaTC",
                "a10d:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hXrAA16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "7564:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhdOn9htOj+hXrAA16FcPhJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC"
            ]
        },
        "rsfirewall": {
            "company": "RSJoomla!",
            "name": "RSFirewall",
            "regex": "COM_RSFIREWALL_",
            "signatures": [
                "d829:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XaTC"
            ]
        },
        "safe3": {
            "company": "Safe3",
            "name": "Safe3",
            "regex": "Server: Safe3 Web Firewall|Safe3waf/",
            "signatures": [
                "1b84:RVZXum60OEhCWKpAYKYPk4JyWOpohM4IiUYMr2RWg1uQJbX2uhdOnthtOj+hX7AB16FcPhJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC"
            ]
        },
        "safedog": {
            "company": "Safedog",
            "name": "Safedog",
            "regex": "Server: Safedog|safedogsite/broswer_logo\\.jpg|404\\.safedog\\.cn/sitedog_stat\\.html|404\\.safedog\\.cn/images/safedogsite/head\\.png",
            "signatures": [
                "0ee1:RVdXu261OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AA16FcPhJOdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "28a0:RVZXu261OEhCWapBYKcPk4JzWOpohM4IiUcMr2RXg1uQJbX3uhdOnthsOj+hX7AA16FcPhJOdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC",
                "90fa:RVZXu261OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AA16FcPhJOdLoXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD"
            ]
        },
        "safeline": {
            "company": "Chaitin Tech",
            "name": "SafeLine Next Gen WAF",
            "regex": "<!\\-\\- event_id: [0-9a-f]{32} \\-\\->",
            "signatures": []
        },
        "secureentry": {
            "company": "United Security Providers",
            "name": "Secure Entry Server",
            "regex": "Server: Secure Entry Server",
            "signatures": [
                "6249:RVZXum60OEhCWKpAYKYPk4JzWOpohM4IiUcMr2RWg1uQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "secureiis": {
            "company": "BeyondTrust",
            "name": "SecureIIS Web Server Security",
            "regex": "//www\\.eeye\\.com/SecureIIS/|\\?subject=[^>]*SecureIIS Error|SecureIIS[^<]+Web Server Protection",
            "signatures": [
                "b43e:RVZXum60OEhCWKpAYKYPkoJzWOtohM4IiUcMrmRWg1qQJbX3uhdOnthsOj+hX7AB16BcPhJOdLoXo2tKaK99n+i6c4VmkwI3FZnxtDtBeq6c36A4chS1XaTC",
                "71c7:RVZXum61OElCWKpAYKYPk4JyWOpohc4IiUYMr2RWg1uQJbX2uhdOnthtOj+hXrAB16FcPhJOdLoXo2tLaK99nui7c4RmkwI2FZjxtDtAeq+c36A4chW1XaTC",
                "f2ed:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJbX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4VmkwI3FZjxtDtAeq6c36A4chS1XaTC"
            ]
        },
        "secupress": {
            "company": "SecuPress",
            "name": "SecuPress",
            "regex": "<h1>SecuPress</h1><h2>\\d{3}",
            "signatures": [
                "bcb4:RVZXum60OEhCWKpAYKYPkoJyWOpohc4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "shieldsecurity": {
            "company": "One Dollar Plugin",
            "name": "Shield Security",
            "regex": "Something in the URL, Form or Cookie data wasn't appropriate",
            "signatures": [
                "e41d:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "389c:RVZXum61OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "a79a:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTD"
            ]
        },
        "securesphere": {
            "company": "Imperva",
            "name": "SecureSphere",
            "regex": "<H2>Error</H2>.+?#FEEE7A.+?<STRONG>Error</STRONG>|Contact support for additional information.<br/>The incident ID is: (\\d{19}|N/A)",
            "signatures": [
                "c055:RVZXum60OEhCWapAYKYPkoJzWOpohM4JiUcMr2RWg1uQJbX2uhZOnthsOj+hX7AB16FcPxJPdLoXomtKaK59n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "f460:RVZXum60OEhCWapBYKYPk4JzWOtohM4JiUcMr2RWg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "9113:RVZXum60OEhCWapBYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "dc2c:RVZXum60OEhCWapBYKYPk4JzWOtohM4JiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "599d:RVZXum60OEhCWapBYKYPk4JzWOtohM4JiUcMr2RWg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "a86e:RVZXum60OEhCWapBYKYPk4JyWOtohM4JiUcMr2RWg1uQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXo2tKaK99n+i6c4RmkgI2FZjxtDtAeq+c36A4chS1XaTC",
                "81ca:RVZXum60OEhCWapBYKYPk4JzWOtohM4IiUcMr2RWg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "siteground": {
            "company": "SiteGround",
            "name": "SiteGround",
            "regex": "The page you are trying to access is restricted due to a security rule|Our system thinks you might be a robot!|/.well-known/captcha/",
            "signatures": [
                "da25:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA5chW1XKTC"
            ]
        },
        "siteguard": {
            "company": "JP-Secure",
            "name": "SiteGuard",
            "regex": "Powered by SiteGuard|The server refuse to browse the page",
            "signatures": [
                "6e49:RVZXum61OElCWapBYKcPk4JzWOtohM4JiUYMr2RWg1qQJbX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "9839:RVZXum61OElCWapBYKcPk4JzWOtohM4JiUYMr2RWg1qQJbX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq6c36A4chS1XaTC",
                "bc2d:RVZXum61OElCWapBYKcPk4JzWOtohM4JiUYMr2RWg1qQJLX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "sitelock": {
            "company": "SiteLock",
            "name": "TrueShield",
            "regex": "SiteLock Incident ID|SiteLock will remember you and will not show this page again|<span class=\\\"value INCIDENT_ID\\\">",
            "signatures": [],
            "note": "Uses Incapsula (Reference: https://www.whitefirdesign.com/blog/2016/11/08/more-evidence-that-sitelocks-trueshield-web-application-firewall-is-really-incapsulas-waf/)"
        },
        "sniper": {
            "company": "Wins",
            "name": "Sniper",
            "regex": "document\\.title = [^;]+Sniper WAF",
            "signatures": []
        },
        "sonicwall": {
            "company": "Dell",
            "name": "SonicWALL",
            "regex": "Server: SonicWALL|(?s)<title>Web Site Blocked</title>.+?nsa_banner",
            "signatures": [
                "f85c:RVZXum61OElCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1qQJLX2uhZOnthsOj+hX7AA16FcPxJPdLoXo2tLaK99nui7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD"
            ]
        },
        "sophos": {
            "company": "Sophos",
            "name": "UTM Web Protection",
            "regex": "Powered by UTM Web Protection",
            "signatures": []
        },
        "squarespace": {
            "company": "Squarespace",
            "name": "Squarespace",
            "regex": "(?s) @ .+?BRICK-50",
            "signatures": [
                "b012:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC",
                "4381:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOn9hsOj6hXrAA16BcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTC"
            ]
        },
        "stackpath": {
            "company": "StackPath",
            "name": "StackPath",
            "regex": "You performed an action that triggered the service and blocked your request",
            "signatures": [
                "5ab0:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUYMr2RWg1uQJbX2uhdOn9hsOj+hXrAA16FcPhJOdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "7e0a:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUYMr2RWg1uQJbX2uhdOn9htOj+hXrAA16FcPxJOdLsXomtKaK59n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD"
            ]
        },
        "sucuri": {
            "company": "Sucuri",
            "name": "Sucuri",
            "regex": "Access Denied - Sucuri Website Firewall|Sucuri WebSite Firewall - CloudProxy - Access Denied|Questions\\?.+cloudproxy@sucuri\\.net",
            "signatures": [
                "60a9:RVZXum61OElCWapAYKYPk4JzWOpohM4JiUYMr2RXg1uQJbX3uhdOn9htOj+hXrAB16FcPxJPdLsXo2tLaK99n+i7c4RmkwI2FZjxtDtAeq+c36A5chW1XaTC"
            ]
        },
        "tencent": {
            "company": "Tencent Cloud Computing",
            "name": "Tencent Cloud|Waterproof Wall",
            "regex": "waf\\.tencent-cloud\\.com|window.location.href=.https://waf.tencent.com/501page.html",
            "signatures": [
                "3f82:RVZXum60OEhCWapBYKcPk4JzWOpohM4IiUYMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tKaK99nui7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTD"
            ]
        },
        "tmg": {
            "company": "Microsoft",
            "name": "Forefront Threat Management Gateway",
            "regex": "",
            "signatures": [
                "4d00:RVZXum60OEhCWKpAYKYPkoJyWOpohM4JiUYMr2RWg1qQJLX3uhdOnthsOj+hX7AB16BcPhJPdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq+c3qA4chS1XaTC"
            ]
        },
        "urlmaster": {
            "company": "iFinity/DotNetNuke",
            "name": "Url Master SecurityCheck",
            "regex": "UrlRewriteModule\\.SecurityCheck|X-UrlMaster-(Debug|Ex):",
            "signatures": [
                "ddd8:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XaTC"
            ]
        },
        "urlscan": {
            "company": "Microsoft",
            "name": "UrlScan",
            "regex": "Rejected-By-UrlScan",
            "signatures": [
                "0294:RVdXum60OEhCWKpAYKYPk4JyWOpohM4IiUYMrmRXg1qQJLX2uhdOn9htOj+hXrAB16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC"
            ]
        },
        "vfw": {
            "company": "OWASP",
            "name": "Varnish Firewall",
            "regex": "Request rejected by xVarnish-WAF",
            "signatures": []
        },
        "virusdie": {
            "company": "Virusdie LLC",
            "name": "Virusdie",
            "regex": "Virusdie</title>|http://cdn\\.virusdie\\.ru/splash/firewallstop\\.png|<meta name=\\\"FW_BLOCK\\\"",
            "signatures": []
        },
        "vsf": {
            "company": "Varnish Cache Project",
            "name": "Varnish Security Firewall",
            "regex": "<title>403 Naughty, not nice!</title>",
            "signatures": [
                "26fa:RVZXum60OEhCWKpAYKYPkoJyWOpohM4JiUcMr2RXg1qQJLX3uhZOnthsOj+hXrAA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTD"
            ]
        },
        "wallarm": {
            "company": "Wallarm",
            "name": "Wallarm",
            "regex": "Server: nginx-wallarm",
            "signatures": [
                "c02b:RVZXu261OElCWapBYKcPk4JzWOpohM4JiUcMr2RWg1uQJbX3uhdOnthsOj+hXrAB16FcPxJOdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC"
            ]
        },
        "wapples": {
            "company": "Penta Security",
            "name": "Wapples",
            "regex": "",
            "signatures": [
                "60b7:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1uQJLX2uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XKTC"
            ]
        },
        "watchguard": {
            "company": "WatchGuard Technologies",
            "name": "WatchGuard",
            "regex": "Server: WatchGuard|Request denied by WatchGuard Firewall",
            "signatures": [
                "4f4f:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RWg1uQJLX2uhZOnthsOj+hXrAA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "2a3c:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RXg1uQJLX2uhZOnthsOj+hX7AA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "aa64:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RXg1uQJLX2uhZOnthsOj+hX7AA16FcPhJOdLoXomtKaK59nui7c4RmkgI3FZjxtDtAeq+c3qA4chW1XaTC"
            ]
        },
        "webarx": {
            "company": "WebARX",
            "name": "WebARX",
            "regex": "/wp-content/plugins/webarx/includes/|This request has been blocked by.+?>WebARX<",
            "signatures": []
        },
        "webknight": {
            "company": "AQTRONIX",
            "name": "WebKnight",
            "regex": "WebKnight Application Firewall Alert|AQTRONIX WebKnight|HTTP Error 999\\.0 - AW Special Error",
            "signatures": [
                "80f9:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJbX2uhdOnthtOj+hXrAB16FcPhJPdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "73e5:RVZXum60OEhCWKpAYKYPk4JyWOtohM4JiUcMrmRXg1uQJbX3uhZOnthsOj6hX7AA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XaTC",
                "d0f0:RVdXum60OEhCWKpAYKYPk4JyWOtohM4JiUcMrmRXg1uQJbX3uhdOn9htOj+hX7AA16FcPxJOdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC",
                "f0c3:RVZXum61OElCWKpAYKYPk4JyWOtohM4JiUcMr2RXg1uQJbX3uhZOnthsOj6hX7AA16BcPhJOdLoXo2tKaK59n+i6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "6763:RVZXum61OElCWKpAYKYPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "7701:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJbX2uhdOn9htOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "902b:RVdXum60OEhCWKpAYKYPk4JyWOpohM4IiUYMrmRXg1qQJbX2uhdOn9htOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c36A4chW1XaTC",
                "4d4d:RVdXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJbX2uhdOn9htOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC",
                "17a8:RVZXum60OEhCWKpAYKYPkoJyWOpohM4JiUcMrmRXg1qQJbX3uhdOnthtOj+hXrAB16FcPhJPdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq+c3qA4chS1XKTC"
            ]
        },
        "webland": {
            "company": "WebLand",
            "name": "WebLand",
            "regex": "Server: Apache Protected by Webland WAF",
            "signatures": [
                "4ba0:RVZXum60OEhCWKpAYKYPkoJzWOpohc4IiUYMr2RWg1uQJLX3uhZOnthsOj6hXrAA16BcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "webseal": {
            "company": "IBM",
            "name": "WebSEAL",
            "regex": "(?i)Server: WebSEAL|This is a WebSEAL error message template file|The Access Manager WebSEAL server received an invalid HTTP request",
            "signatures": [
                "0338:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRWg1qQJLX2uhZOnthtOj+hXrAA16FcPhJOdLoXomtKaK59nui6c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC"
            ]
        },
        "webtotem": {
            "company": "WebTotem",
            "name": "WebTotem",
            "regex": "The current request was blocked by.+?>WebTotem<",
            "signatures": []
        },
        "wordfence": {
            "company": "Defiant",
            "name": "Wordfence",
            "regex": "Generated by Wordfence|This response was generated by Wordfence|broke one of the Wordfence (advanced )?blocking rules|: wfWAF|/plugins/wordfence",
            "signatures": [
                "d04a:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "26b1:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAA16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "09cf:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtBeq6c3qA4chW1XaTC",
                "1834:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RXg1uQJLX3uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c36A4chW1XaTC",
                "d38c:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkwI3FZjxtDtAeq6c3qA4chW1XaTC",
                "d5bb:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1uQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "3f1c:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTD",
                "dbfe:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA5chW1XaTC",
                "5b85:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMr2RXg1uQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA5chW1XaTD",
                "f806:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hX7AB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "0f0d:RVZXum61OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkwI3FZjxtDtAeq6c3qA4chW1XaTC",
                "b13e:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJbX3uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "40eb:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16BcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XaTC",
                "93cd:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chS1XKTC",
                "ba7d:RVZXum60OEhCWKpAYKYPkoJyWOpohM4IiUYMrmRXg1qQJLX2uhdOnthtOj+hXrAB16FcPxJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq6c3qA4chW1XKTC"
            ]
        },
        "wts": {
            "company": "WTS",
            "name": "WTS",
            "regex": "Server: wts/|>WTS\\-WAF",
            "signatures": [
                "e94f:RVZXum61OElCWapAYKYPkoJzWOpohM4JiUcMr2RXg1uQJLX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XKTC",
                "12ce:RVZXum61OElCWapAYKYPkoJzWOpohM4IiUYMr2RWg1uQJLX3uhdOnthtOj+hX7AB16FcPhJPdLsXo2tKaK99n+i7c4RmkgI2FZjxtDtAeq+c3qA4chW1XKTC"
            ]
        },
        "yundun": {
            "company": "Yundun",
            "name": "Yundun",
            "regex": "Blocked by YUNDUN Cloud WAF|yundun\\.com/yd_http_error/",
            "signatures": [
                "4853:RVZXum61OEhCWapBYKcPk4JzWOtohM4JiUcMr2RXg1uQJbX3uhdOnthtOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4RmkgI2FZjxtDtAeq+c36A5chW1XaTC"
            ]
        },
        "yunsuo": {
            "company": "Yunsuo",
            "name": "Yunsuo",
            "regex": "yunsuo_session|<img class=\\\"yunsuologo\\\"",
            "signatures": [
                "441b:RVZXum60OEhCWKpAYKYPkoJzWOtohM4JiUcMr2RXg1uQJbX3uhdOnthsOj+hX7AA16FcPxJOdLoXomtKaK59nui7c4VmkgI2FZjxtDtAeq+c3qA4chW1XKTC",
                "e795:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJbX3uhdOnthsOj+hX7AB16FcPhJPdLsXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XaTC",
                "7b8e:RVZXum60OEhCWKpAYKYPkoJzWOpohM4JiUcMr2RXg1uQJbX3uhdOnthsOj+hX7AA16FcPhJOdLoXomtKaK59nui7c4RmkgI2FZjxtDtAeq+c3qA4chW1XKTC"
            ]
        },
        "zenedge": {
            "company": "Zenedge",
            "name": "Zenedge",
            "regex": "(?s)Server: ZENEDGE.+?<div class=\\\"number\\\">403</div>",
            "signatures": [
                "a8fb:RVdXu260OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI2FZnxtDtBeq+c36A4chW1XaTD",
                "ba3d:RVdXu260OEhCWapBYKcPk4JzWOpohM4JiUcMr2RXg1uQJbX3uhdOn9htOj+hX7AB16FcPxJPdLsXo2tLaK99n+i7c4VmkwI2FZjxtDtAeq+c36A4chW1XaTD"
            ]
        }
    }
}
`
