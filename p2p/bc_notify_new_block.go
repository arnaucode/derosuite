// Copyright 2017-2018 DERO Project. All rights reserved.
// Use of this source code in any form is governed by RESEARCH license.
// license can be found in the LICENSE file.
// GPG: 0F39 E425 8C65 3947 702A  8234 08B2 0360 A03A 9DE8
//
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
// STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF
// THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package p2p

import "bytes"
import "encoding/binary"

import "github.com/romana/rlog"

import "github.com/arnaucode/derosuite/block"
import "github.com/arnaucode/derosuite/transaction"

// FIXME this code can also be shared by NOTIFY_NEW_BLOCK, NOTIFY_NEW_TRANSACTIONS, Handle_BC_Notify_Response_GetObjects
// this code handles a new block floating in the network
func Handle_BC_Notify_New_Block(connection *Connection,
	i_command_header *Levin_Header, buf []byte) {
	var bl block.Block
	var cbl block.Complete_Block

	connection.logger.Debugf("Incoming NOTIFY_NEW_BLOCK")

	// deserialize data header
	var i_data_header Levin_Data_Header // incoming data header

	err := i_data_header.DeSerialize(buf)

	if err != nil {
		connection.logger.Debugf("We should destroy connection here, data header cnot deserialized")
		return
	}

	buf = i_data_header.Data

	pos := bytes.Index(buf, []byte("block"))
	// find inner position of block
	pos = pos + 6 // jump to varint length position and decode

	buf = buf[pos:]
	block_length, done := Decode_Boost_Varint(buf)
	rlog.Tracef(2, "Block length %d %x\n", block_length, buf[:8])
	buf = buf[done:]

	block_buf := buf[:block_length]

	err = bl.Deserialize(block_buf)
	if err != nil {
		connection.logger.Warnf("Block could not be deserialized successfully err %s\n", err)
		connection.logger.Debugf("We should destroy connection here, block not deserialized")
		return
	}

	hash := bl.GetHash()
	rlog.Tracef(1, "Block deserialized successfully  %x\n", hash[:32])
	rlog.Tracef(1, "Tx hash length %d\n", len(bl.Tx_hashes))
	for i := range bl.Tx_hashes {
		rlog.Tracef(2, "%d tx %x\n", i, bl.Tx_hashes[i][:32])
	}
	// point buffer to check whether any more tx exist
	buf = buf[block_length:]

	pos = bytes.Index(buf, []byte("\x03txs\x8a")) // at this point to

	if pos > -1 {
		rlog.Tracef(3, "txt pos %d\n", pos)

		buf = buf[pos+5:]
		// decode remain data length ( though we know it from buffer size, but still verify it )

		tx_count, done := Decode_Boost_Varint(buf)
		buf = buf[done:]

		for i := uint64(0); i < tx_count; i++ {

			var tx transaction.Transaction

			tx_len, done := Decode_Boost_Varint(buf)
			buf = buf[done:]
			rlog.Tracef(3, "tx count %d  i %d  tx_len %d\n", tx_count, i, tx_len)

			tx_bytes := buf[:tx_len]

			// deserialize and verrify transaction

			err = tx.DeserializeHeader(tx_bytes)

			if err != nil {
				connection.logger.Warnf("Transaction could not be deserialized\n")

			} else {
				hash := tx.GetHash()
				rlog.Tracef(2, "Transaction deserialised successfully  hash %x\n", hash[:32])

				// add tx to block chain, we must verify that the tx has been mined
				// add all transaction to TX pool , if not added
				//chain.Add_TX(&tx)
				cbl.Txs = append(cbl.Txs, &tx)
			}

			buf = buf[tx_len:] // setup for next tx

		}
	}

	height_string := []byte("\x19current_blockchain_height\x05")
	pos = bytes.Index(buf, height_string) // at this point to

	if pos < 0 {
		connection.logger.Debugf("We should destroy connection here, block not deserialized")
		return
	}

	pos = pos + len(height_string)
	buf = buf[pos:]
	current_peer_height := binary.LittleEndian.Uint64(buf)

	connection.Last_Height = current_peer_height

	//connection.logger.Infof("buffer height %x  current height %d   complete %d\n", buf, connection.Last_Height, complete_block)

	// at this point, if it's a block we should try to add it to block chain
	// try to add block to chain
	connection.logger.Debugf("Found new  block adding it to chain %s", bl.GetHash())

	// TODO check returned status, either drop connection or replay
	cbl.Bl = &bl
	chain.Add_Complete_Block(&cbl)

}
