/*
Copyright © 2026 irrisdev lithium8260@proton.me

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package internal

const (
	HeaderSize      = 20        // bytes
	ChunkHeaderSize = 16        // bytes
	RWSize          = 64 * 1024 // 64 KB
	ChunkSize       = 1 << 20   // 1 MiB
)

var MagicHeader = [4]byte{'g', 'e', 'n', 'c'}

type Header struct {
	Magic [4]byte
	Salt  [16]byte
}

type ChunkHeader struct {
	Nonce  [12]byte
	Length uint32
}
