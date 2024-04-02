package otime

import "strings"

var NameLocaleID = strings.NewReplacer(
	"Sunday", "Minggu",
	"Monday", "Senin",
	"Tuesday", "Selasa",
	"Wednesday", "Rabu",
	"Thursday", "Kamis",
	"Friday", "Jumat",
	"Saturday", "Sabtu",
	"January", "Januari",
	"February", "Februari",
	"March", "Maret",
	"April", "April",
	"May", "Mei",
	"June", "Juni",
	"July", "Juli",
	"August", "Agustus",
	"September", "September",
	"October", "Oktober",
	"November", "November",
	"December", "Desember",
)
