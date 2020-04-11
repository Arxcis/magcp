package main

import (
	"fmt"
	"regexp"
	"net/url"
)

type Magnet struct {
	xt string
	dn string
	tr []string
}

var reg_xt *regexp.Regexp = regexp.MustCompile(`xt=([^&]+)`)
var reg_dn *regexp.Regexp = regexp.MustCompile(`dn=([^&]+)`)
var reg_tr *regexp.Regexp = regexp.MustCompile(`tr=([^&]+)`)

func UnmarshalMagnet(in_magnet string, out_magnet *Magnet) {
	// xt
	match_xt := reg_xt.FindStringSubmatch(in_magnet)
	out_magnet.xt = match_xt[1]

	// dn
	match_dn := reg_dn.FindStringSubmatch(in_magnet)
	matchstr, _ := url.QueryUnescape(match_dn[1])
	out_magnet.dn = matchstr

	// tr
	match_tr := reg_tr.FindAllStringSubmatch(in_magnet, -1)
	for _, match := range match_tr {
		matchstr, _ = url.QueryUnescape(match[1])
		out_magnet.tr = append(out_magnet.tr, matchstr)
	}
}

func main() {
	const magnet1str = `magnet:?
xt=urn:btih:89599BF4DC369A3A8ECA26411C5CCF922D78B486&
dn=Interstellar+%282014%29+1080p+BrRip+x264+-+YIFY&
tr=udp%3A%2F%2Ftracker.yify-torrents.com%2Fannounce&
tr=udp%3A%2F%2Fopen.demonii.com%3A1337&
tr=udp%3A%2F%2Fexodus.desync.com%3A6969&
tr=udp%3A%2F%2Ftracker.istole.it%3A80&
tr=udp%3A%2F%2Ftracker.publicbt.com%3A80&
tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&
tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&
tr=udp%3A%2F%2F9.rarbg.com%3A2710&
tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&
tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&
tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&
tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce`

	const magnet2str = `magnet:?
xt=urn:btih:098AFE45F8BE4138BBA50A746B1B580DB17B1696&
dn=Cosmos+%282019%29+%5BWEBRip%5D+%5B720p%5D+%5BYTS%5D+%5BYIFY%5D&
tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&
tr=udp%3A%2F%2F9.rarbg.com%3A2710%2Fannounce&
tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&
tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337&
tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&
tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&
tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&
tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce`

	magnet1 := Magnet{}
	magnet2 := Magnet{}

	// Unmarshal magnet1 string
	{
		fmt.Println()

		Unmarshal(magnet1str, &magnet1)
		Unmarshal(magnet2str, &magnet2)

		fmt.Printf("%+#v", magnet1)
		fmt.Println()
		fmt.Println()
		fmt.Printf("%+#v", magnet2)
		fmt.Println()
	}
}

