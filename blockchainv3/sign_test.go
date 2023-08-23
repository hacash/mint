package blockchainv3

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"testing"
)

func Test_trs_sign(t *testing.T) {

	addr1, _ := fields.CheckReadableAddress("1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y")
	addr2, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")

	amt1 := fields.NewAmountNumSmallCoin(1)

	act1 := actions.NewAction_1_SimpleToTransfer(*addr2, amt1)

	trs1, _ := transactions.NewEmptyTransaction_2_Simple( *addr1 )
	trs1.AppendAction( act1 )

	fmt.Println( trs1 )


	//--------------------------------

	blkdata1, _ := hex.DecodeString("01000005c8d2006320a58800000000007cb0615db42c344fccd01d6da05473aeadc205e6aa225d43117eb3bc02b9b85ca5bf768397de4ae9bd266f51aab2ce931ac2128978731341c13f4600000005eaf2def5d8afc93a000000001a9ffb585dd247d351dae9c18031d7e3c762d7a2f80103576f575f3035323700000000000000c301000000000000000000000000000000000000000000000000000000000000000000020063209b1f004b7812e33a7774757269f4dad1abf85aa422d340f60204c700010004485341594b4d00e4e8000000000040b6cb35caaedd1776a8a24e7dd7f918598d9a6912b1498d73fc7a000000008f0acb00004b7812e33a7774757269f4dad1abf85aa422d340d612d69e7c1a9901c4ec2c6f8249eb3423779b615b06e8ba1f581df172d85bfd000103f76807b92ff7afba6d436e67fc1e2ed0817b61d00dd700543a5a6fc04835deaede18fb0c4e1925ff120f63962ca3a4ede2d754e59bbf5f38b441059cbdf35a7b353745d429c2acbf2e923638297de9bb8a44af07ac0194f8c823eff10979d1fd000002006320a2ec00fc6038915c31971586e3e1a7b9d49e0baa7621a9f601010001000100fe4db49e64c250009abdb403c27bcd323ef93539f8023b3100010247cb19cdbcf94c513a9c2a6feb4e8a8da08c0165eedb55ef61732307fd9a2f5203ab0a110108b1a823e3d22a1e55bbe570f11fc25b1bf7a4c2fc5eddee37f8ca044af015411720727aa46fd49bd2a69073a220c6d74014e4090a99c5d9bd483c000002006320a365006803f6f526aa2a85aa39142c14aa12ba88ba8867f6010100010001003c75683c14e7b5948f68c8b6501c6a6e919a3708f80207df00010249a10aca1b2d87066d6f6465bd99cc564aa1be60d3d9f82d08194df604312fa1830f53f68504598a7a993c602c2e1bc904b0e5fbfaab92305c6c7cf3886541e419da187131d93946fd0cd35932b28128d3798c6a4bbac409be3028f57f3836b9000002006320a4c1007a006db6b69c1371936b4eceec674d5c71a9d8c6f6010100010006007a006db6b69c1371936b4eceec674d5c71a9d8c600d60d72ec831078d7bdbe145418920c18fd1580040e494853425a4d41484d45425559564e454d4d494154425748544d414b575553415749495454455a425453534958565441564b454d4854544841484958544e4d4149594b4e574541494549424948494945554b4e5400010230c3beca7eba3aa8817b2ce00686a0408df58a1012cca22467960711ab39a5a93c17eade72ea7f244d5f1324f0fa571f552a03f9854740b0a8e1e20ca99fca436553fdfe20426d4cbfae1752eb3e043e19bbb6cc1baa557541d7d195ebe01c840000")

	blk, _, _ := blocks.ParseBlock(blkdata1, 0)


	trslist := blk.GetTrsList()
	tartrs := trslist[2]
	fmt.Println( tartrs.GetAddress().ToReadable() )

	trsedit := tartrs.(*transactions.Transaction_2_Simple)
	//trsedit.CleanSigns()
	//trsedit.Signs[0].PublicKey[1] = 8
	trsedit.Signs[0].Signature[1] = 8 // change signdata


	//fmt.Println( trsedit.VerifyAllNeedSigns() )
	fmt.Println( blk.VerifyNeedSigns() )






}
