// Copyright (c) - Damien Fontaine <damien.fontaine@lineolia.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package utils

import (
	"bytes"
	"regexp"
	"strings"
)

//Liste des accents
var a = map[rune]string{
	'À': "A",
	'Á': "A",
	'Â': "A",
	'Ã': "A",
	'Ä': "A",
	'Æ': "AE",
	'Ç': "C",
	'È': "E",
	'É': "E",
	'Ê': "E",
	'Ë': "E",
	'Ì': "I",
	'Í': "I",
	'Î': "I",
	'Ï': "I",
	'Ñ': "N",
	'Ò': "O",
	'Ó': "O",
	'Ô': "O",
	'Õ': "O",
	'Ö': "O",
	'Œ': "OE",
	'Ù': "U",
	'Ú': "U",
	'Ü': "U",
	'Û': "U",
	'Ý': "Y",
	'à': "a",
	'á': "a",
	'â': "a",
	'ã': "a",
	'ä': "a",
	'æ': "ae",
	'ç': "c",
	'è': "e",
	'é': "e",
	'ê': "e",
	'ë': "e",
	'ì': "i",
	'í': "i",
	'î': "i",
	'ï': "i",
	'ñ': "n",
	'ń': "n",
	'ò': "o",
	'ó': "o",
	'ô': "o",
	'õ': "o",
	'ō': "o",
	'ö': "o",
	'œ': "oe",
	'ś': "s",
	'ù': "u",
	'ú': "u",
	'û': "u",
	'ū': "u",
	'ü': "u",
	'ý': "y",
	'ÿ': "y",
	'ż': "z",
}

//SanitizeAccent remplace les caractères accentués
func SanitizeAccent(s string) string {
	b := bytes.NewBufferString("")
	for _, c := range s {
		if v, p := a[c]; p {
			b.WriteString(v)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

//SanitizeTitle transforme un titre en chaîne de caractères utilisable pour l'affichage dans une URL
func SanitizeTitle(s string) string {
	p := SanitizeAccent(s)
	r, _ := regexp.Compile("[^A-Za-z0-9 ']+")
	p = r.ReplaceAllString(p, "")
	p = strings.Replace(p, " ", "-", -1)
	r, _ = regexp.Compile("^-|-$")
	p = r.ReplaceAllString(p, "")
	p = strings.Replace(p, "'", "-", -1)
	r, _ = regexp.Compile("-+")
	p = r.ReplaceAllString(p, "-")
	p = strings.Replace(p, "&", "and", -1)
	p = strings.ToLower(p)
	return p
}
