////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: context_test.go
// @Date: 2018/08/08 16:49:08
////////////////////////////////////////////////////////////////////////////////

package actioncontext_test

import (
	"bchain.io/common/types"
	"bchain.io/core/actioncontext"
	"bchain.io/core/interpreter"
	"bchain.io/core/interpreter/jsre"
	"bchain.io/core/state"
	"bchain.io/core/transaction"
	"bchain.io/utils/crypto"
	"bchain.io/utils/database"
	"strings"
	"testing"
	"time"
)

func TestContext_1(t *testing.T) {
	interpreter.Singleton().Register(func() interpreter.PluginImpl {
		return &jsre.JSRE{}
	})

	interpreter.Singleton().Initialize()
	interpreter.Singleton().Startup()

	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()
	for i := 0; i < 10; i++ {
		go func() {
			tr := transaction.Action{}
			blkctx := actioncontext.NewBlockContext(stateDb, db, tmpdb, nil, types.Address{})
			ctx := actioncontext.NewContext(types.Address{}, &tr, blkctx)

			ctx.InitForTest()
			ctx.Exec(interpreter.Singleton())
		}()
	}

	time.Sleep(300000 * time.Microsecond)
	interpreter.Singleton().Shutdown()
	time.Sleep(1 * time.Second)
}

func TestContextApiSha256(t *testing.T) {
	ctx := actioncontext.Context{}
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(&ctx)

	preImage := []byte("hello world\n")
	exp := "0xa948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"
	aim := ctxApi.Crypto.Sha256(preImage)
	if exp != types.ToHex(aim[:]) {
		t.Errorf("ActionContext Crypto.Sha256 = %s; want %s", exp, aim)
	}

}

func TestContextApiSha1(t *testing.T) {
	type sha1Test struct {
		out string
		in  string
	}

	var golden = []sha1Test{
		{"0x76245dbf96f661bd221046197ab8b9f063f11bad", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n"},
		{"0xda39a3ee5e6b4b0d3255bfef95601890afd80709", ""},
		{"0x86f7e437faa5a7fce15d1ddcb9eaeaea377667b8", "a"},
		{"0xda23614e02469a0d7c7bd1bdab5c9c474b1904dc", "ab"},
		{"0xa9993e364706816aba3e25717850c26c9cd0d89d", "abc"},
		{"0x81fe8bfe87576c3ecb22426f8e57847382917acf", "abcd"},
		{"0x03de6c570bfe24bfc328ccd7ca46b76eadaf4334", "abcde"},
		{"0x1f8ac10f23c5b5bc1167bda84b833e5c057a77d2", "abcdef"},
		{"0x2fb5e13419fc89246865e7a324f476ec624e8740", "abcdefg"},
		{"0x425af12a0743502b322e93a015bcf868e324d56a", "abcdefgh"},
		{"0xc63b19f1e4c8b5f76b25c49b8b87f57d8e4872a1", "abcdefghi"},
		{"0xd68c19a0a345b7eab78d5e11e991c026ec60db63", "abcdefghij"},
		{"0xebf81ddcbe5bf13aaabdc4d65354fdf2044f38a7", "Discard medicine more than two years old."},
		{"0xe5dea09392dd886ca63531aaa00571dc07554bb6", "He who has a shady past knows that nice guys finish last."},
		{"0x45988f7234467b94e3e9494434c96ee3609d8f8f", "I wouldn't marry him with a ten foot pole."},
		{"0x55dee037eb7460d5a692d1ce11330b260e40c988", "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
		{"0xb7bc5fb91080c7de6b582ea281f8a396d7c0aee8", "The days of the digital watch are numbered.  -Tom Stoppard"},
		{"0xc3aed9358f7c77f523afe86135f06b95b3999797", "Nepal premier won't resign."},
		{"0x6e29d302bf6e3a5e4305ff318d983197d6906bb9", "For every action there is an equal and opposite government program."},
		{"0x597f6a540010f94c15d71806a99a2c8710e747bd", "His money is twice tainted: 'taint yours and 'taint mine."},
		{"0x6859733b2590a8a091cecf50086febc5ceef1e80", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
		{"0x514b2630ec089b8aee18795fc0cf1f4860cdacad", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
		{"0xc5ca0d4a7b6676fc7aa72caa41cc3d5df567ed69", "size:  a.out:  bad magic"},
		{"0x74c51fa9a04eadc8c1bbeaa7fc442f834b90a00a", "The major problem is with sendmail.  -Mark Horton"},
		{"0x0b4c4ce5f52c3ad2821852a8dc00217fa18b8b66", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
		{"0x3ae7937dd790315beb0f48330e8642237c61550a", "If the enemy is within range, then so are you."},
		{"0x410a2b296df92b9a47412b13281df8f830a9f44b", "It's well we cannot hear the screams/That we create in others' dreams."},
		{"0x841e7c85ca1adcddbdd0187f1289acb5c642f7f5", "You remind me of a TV show, but that's all right: I watch it anyway."},
		{"0x163173b825d03b952601376b25212df66763e1db", "C is as portable as Stonehedge!!"},
		{"0x32b0377f2687eb88e22106f133c586ab314d5279", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
		{"0x0885aaf99b569542fd165fa44e322718f4a984e0", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
		{"0x6627d6904d71420b0bf3886ab629623538689f45", "How can you write a big system without C++?  -Paul Glick"},
	}

	ctx := actioncontext.Context{}
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(&ctx)

	for _, item := range golden {
		exp := item.out
		aim := ctxApi.Crypto.Sha1([]byte(item.in))
		if exp != types.ToHex(aim[:]) {
			t.Errorf("ActionContext Crypto.Sha512 = %s; want %s", exp, aim)
		}
	}
}

func TestContextApiSha512(t *testing.T) {
	type sha512Test struct {
		in     string
		out512 string
	}

	var golden = []sha512Test{
		{
			"",
			"0xcf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			"a",
			"0x1f40fc92da241694750979ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
		},
		{
			"ab",
			"0x2d408a0717ec188158278a796c689044361dc6fdde28d6f04973b80896e1823975cdbf12eb63f9e0591328ee235d80e9b5bf1aa6a44f4617ff3caf6400eb172d",
		},
		{
			"abc",
			"0xddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
		},
		{
			"abcd",
			"0xd8022f2060ad6efd297ab73dcc5355c9b214054b0d1776a136a669d26a7d3b14f73aa0d0ebff19ee333368f0164b6419a96da49e3e481753e7e96b716bdccb6f",
		},
		{
			"abcde",
			"0x878ae65a92e86cac011a570d4c30a7eaec442b85ce8eca0c2952b5e3cc0628c2e79d889ad4d5c7c626986d452dd86374b6ffaa7cd8b67665bef2289a5c70b0a1",
		},
		{
			"abcdef",
			"0xe32ef19623e8ed9d267f657a81944b3d07adbb768518068e88435745564e8d4150a0a703be2a7d88b61e3d390c2bb97e2d4c311fdc69d6b1267f05f59aa920e7",
		},
		{
			"abcdefg",
			"0xd716a4188569b68ab1b6dfac178e570114cdf0ea3a1cc0e31486c3e41241bc6a76424e8c37ab26f096fc85ef9886c8cb634187f4fddff645fb099f1ff54c6b8c",
		},
		{
			"abcdefgh",
			"0xa3a8c81bc97c2560010d7389bc88aac974a104e0e2381220c6e084c4dccd1d2d17d4f86db31c2a851dc80e6681d74733c55dcd03dd96f6062cdda12a291ae6ce",
		},
		{
			"abcdefghi",
			"0xf22d51d25292ca1d0f68f69aedc7897019308cc9db46efb75a03dd494fc7f126c010e8ade6a00a0c1a5f1b75d81e0ed5a93ce98dc9b833db7839247b1d9c24fe",
		},
		{
			"abcdefghij",
			"0xef6b97321f34b1fea2169a7db9e1960b471aa13302a988087357c520be957ca119c3ba68e6b4982c019ec89de3865ccf6a3cda1fe11e59f98d99f1502c8b9745",
		},
		{
			"Discard medicine more than two years old.",
			"0x2210d99af9c8bdecda1b4beff822136753d8342505ddce37f1314e2cdbb488c6016bdaa9bd2ffa513dd5de2e4b50f031393d8ab61f773b0e0130d7381e0f8a1d",
		},
		{
			"He who has a shady past knows that nice guys finish last.",
			"0xa687a8985b4d8d0a24f115fe272255c6afaf3909225838546159c1ed685c211a203796ae8ecc4c81a5b6315919b3a64f10713da07e341fcdbb08541bf03066ce",
		},
		{
			"I wouldn't marry him with a ten foot pole.",
			"0x8ddb0392e818b7d585ab22769a50df660d9f6d559cca3afc5691b8ca91b8451374e42bcdabd64589ed7c91d85f626596228a5c8572677eb98bc6b624befb7af8",
		},
		{
			"Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave",
			"0x26ed8f6ca7f8d44b6a8a54ae39640fa8ad5c673f70ee9ce074ba4ef0d483eea00bab2f61d8695d6b34df9c6c48ae36246362200ed820448bdc03a720366a87c6",
		},
		{
			"The days of the digital watch are numbered.  -Tom Stoppard",
			"0xe5a14bf044be69615aade89afcf1ab0389d5fc302a884d403579d1386a2400c089b0dbb387ed0f463f9ee342f8244d5a38cfbc0e819da9529fbff78368c9a982",
		},
		{
			"Nepal premier won't resign.",
			"0x420a1faa48919e14651bed45725abe0f7a58e0f099424c4e5a49194946e38b46c1f8034b18ef169b2e31050d1648e0b982386595f7df47da4b6fd18e55333015",
		},
		{
			"For every action there is an equal and opposite government program.",
			"0xd926a863beadb20134db07683535c72007b0e695045876254f341ddcccde132a908c5af57baa6a6a9c63e6649bba0c213dc05fadcf9abccea09f23dcfb637fbe",
		},
		{
			"His money is twice tainted: 'taint yours and 'taint mine.",
			"0x9a98dd9bb67d0da7bf83da5313dff4fd60a4bac0094f1b05633690ffa7f6d61de9a1d4f8617937d560833a9aaa9ccafe3fd24db418d0e728833545cadd3ad92d",
		},
		{
			"There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977",
			"0xd7fde2d2351efade52f4211d3746a0780a26eec3df9b2ed575368a8a1c09ec452402293a8ea4eceb5a4f60064ea29b13cdd86918cd7a4faf366160b009804107",
		},
		{
			"It's a tiny change to the code and not completely disgusting. - Bob Manchek",
			"0xb0f35ffa2697359c33a56f5c0cf715c7aeed96da9905ca2698acadb08fbc9e669bf566b6bd5d61a3e86dc22999bcc9f2224e33d1d4f32a228cf9d0349e2db518",
		},
		{
			"size:  a.out:  bad magic",
			"0x3d2e5f91778c9e66f7e061293aaa8a8fc742dd3b2e4f483772464b1144189b49273e610e5cccd7a81a19ca1fa70f16b10f1a100a4d8c1372336be8484c64b311",
		},
		{
			"The major problem is with sendmail.  -Mark Horton",
			"0xb2f68ff58ac015efb1c94c908b0d8c2bf06f491e4de8e6302c49016f7f8a33eac3e959856c7fddbc464de618701338a4b46f76dbfaf9a1e5262b5f40639771c7",
		},
		{
			"Give me a rock, paper and scissors and I will move the world.  CCFestoon",
			"0xd8c92db5fdf52cf8215e4df3b4909d29203ff4d00e9ad0b64a6a4e04dec5e74f62e7c35c7fb881bd5de95442123df8f57a489b0ae616bd326f84d10021121c57",
		},
		{
			"If the enemy is within range, then so are you.",
			"0x19a9f8dc0a233e464e8566ad3ca9b91e459a7b8c4780985b015776e1bf239a19bc233d0556343e2b0a9bc220900b4ebf4f8bdf89ff8efeaf79602d6849e6f72e",
		},
		{
			"It's well we cannot hear the screams/That we create in others' dreams.",
			"0x00b4c41f307bde87301cdc5b5ab1ae9a592e8ecbb2021dd7bc4b34e2ace60741cc362560bec566ba35178595a91932b8d5357e2c9cec92d393b0fa7831852476",
		},
		{
			"You remind me of a TV show, but that's all right: I watch it anyway.",
			"0x91eccc3d5375fd026e4d6787874b1dce201cecd8a27dbded5065728cb2d09c58a3d467bb1faf353bf7ba567e005245d5321b55bc344f7c07b91cb6f26c959be7",
		},
		{
			"C is as portable as Stonehedge!!",
			"0xfabbbe22180f1f137cfdc9556d2570e775d1ae02a597ded43a72a40f9b485d500043b7be128fb9fcd982b83159a0d99aa855a9e7cc4240c00dc01a9bdf8218d7",
		},
		{
			"Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley",
			"0x2ecdec235c1fa4fc2a154d8fba1dddb8a72a1ad73838b51d792331d143f8b96a9f6fcb0f34d7caa351fe6d88771c4f105040e0392f06e0621689d33b2f3ba92e",
		},
		{
			"The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule",
			"0x7ad681f6f96f82f7abfa7ecc0334e8fa16d3dc1cdc45b60b7af43fe4075d2357c0c1d60e98350f1afb1f2fe7a4d7cd2ad55b88e458e06b73c40b437331f5dab4",
		},
		{
			"How can you write a big system without C++?  -Paul Glick",
			"0x833f9248ab4a3b9e5131f745fda1ffd2dd435b30e965957e78291c7ab73605fd1912b0794e5c233ab0a12d205a39778d19b83515d6a47003f19cdee51d98c7e0",
		},
	}

	ctx := actioncontext.Context{}
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(&ctx)

	for _, item := range golden {
		exp := item.out512
		aim := ctxApi.Crypto.Sha512([]byte(item.in))
		if exp != types.ToHex(aim[:]) {
			t.Errorf("ActionContext Crypto.Sha512 = %s; want %s", exp, aim)
		}
	}
}

func TestContextApiRecover(t *testing.T) {
	ctx := actioncontext.Context{}
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(&ctx)

	var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
	var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

	key, _ := crypto.HexToECDSA(testPrivHex)
	addr := types.HexToAddress(testAddrHex)

	msg := crypto.Keccak256([]byte("foo"))
	sig, err := crypto.Sign(msg, key)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	recoveredPub, err := ctxApi.Crypto.Recover(msg, sig)
	if err != nil {
		t.Errorf("ctxApi.Crypto.Recover error: %s", err)
	}
	if string(addr[:]) != string(recoveredPub) {
		t.Errorf("Address mismatch: want: %x have: %x", addr, recoveredPub)
	}

}

func TestContextApiAuth(t *testing.T) {
	type authtest struct {
		addr       string
		isAddr     bool
		isAccount  bool
		isContract bool
		code       string
		key        string
		value      string
	}

	var cases = []authtest{
		{"", false, false, false, "", "", ""},
		{"address", false, false, false, "", "", ""},
		{"970e8128ab834e8eac17ab8e3812f010678cf791999", false, false, false, "", "", ""},
		{"970e8128ab834e8ea!@#$%^&*(12f010678cf79x", false, false, false, "", "", ""},
		{"970e8128ab834e8eac17ab8e3812f010678cf791", true, false, false, "", "", ""},
		{"970e8128ab834e8eac17ab8e3812f010678c970e", true, true, false, "", "this is key", "this is value"},
		{"970e8128ab834e8eac17ab8e3812f010678c8128", true, false, true, "this is code for contract", "", ""},
	}
	// Create an empty state database
	db, _ := database.OpenMemDB()
	state, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()
	for i, v := range cases {
		if v.isAddr {
			addr := types.HexToAddress(v.addr)
			if v.isAccount {
				state.SetNonce(addr, uint64(i+1))
			}
			if v.isContract {
				state.SetCode(addr, []byte(v.code))
			} else if v.isAccount {
				state.SetState(addr, types.BytesToHash([]byte(v.key)), types.BytesToHash([]byte(v.value)))
			}
		}
	}

	tr := transaction.Action{}
	blkctx := actioncontext.NewBlockContext(state, db, tmpdb, nil, types.Address{})
	ctx := actioncontext.NewContext(types.Address{}, &tr, blkctx)
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(ctx)

	for _, v := range cases {
		isAddr := ctxApi.Auth.IsHexAddress(v.addr)
		if isAddr != v.isAddr {
			t.Errorf("TestContextApiAuth addr: %s is %v ;want %v", v.addr, isAddr, v.isAddr)
		}
		if !v.isAddr {
			continue
		}

		isAccount := ctxApi.Auth.IsAccount(v.addr)
		if isAccount != v.isAccount {
			t.Errorf("TestContextApiAuth addr: %s is account %v ;want %v", v.addr, isAccount, v.isAccount)
		}
		if v.isAccount {
			continue
		}

		isContract := ctxApi.Auth.IsContract(v.addr)
		if isContract != v.isContract {
			t.Errorf("TestContextApiAuth addr: %s is contract %v ;want %v", v.addr, isContract, v.isContract)
		}

	}
}

func TestContextApiContract(t *testing.T) {
	type test_st struct {
		addr        string
		interpreter string
		creator     string
		code        string
	}

	var cases = []test_st{
		{
			"0xc57b24b70f335db6d363e47031a61a653c6ffd69",
			"0xcfde91645b17a17d158271ae73a0c18b803d41d2072cdef747e8bbf135e6fe6b",
			"0xffff8128ab834e8eac17ab8e3812f010678cf791",
			`{"inter_name":"jsre.JSRE","code":"dmFyIGN0eGFwaSA9IGN0eEFwaSgpOwoKZnVuY3Rpb24gZ2V0TmFtZSgpewogICAgcmV0dXJuICJ0ZXN0Igp9CgpmdW5jdGlvbiBnZXRWZXJzaW9uKCl7CiAgICByZXR1cm4gIlYwLjAuMSIKfQoKOyhmdW5jdGlvbiAoKSB7CgpjdHhhcGkuY29uc29sZS5wcmludCgiPT09PT09PT09PUNyZWF0ZUNvbnRyYWN0PT09PT09IikKCn0pKHRoaXMpOwoK"}`,
		},
	}
	wantedCode := `var ctxapi = ctxApi();

function getName(){
    return "test"
}

function getVersion(){
    return "V0.0.1"
}

;(function () {

ctxapi.console.print("==========CreateContract======")

})(this);

`
	// Create an empty state database
	db, _ := database.OpenMemDB()
	stateDb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	tmpdb, _ := database.OpenMemDB()

	tr := transaction.Action{}
	blkctx := actioncontext.NewBlockContext(stateDb, db, tmpdb, nil, types.Address{})
	ctx := actioncontext.NewContext(types.Address{}, &tr, blkctx)
	ctx.InitForTest()
	ctxApi := actioncontext.NewAPIs(ctx)

	for _, v := range cases {
		addr := ctxApi.Contract.Create(v.creator, v.code)
		strAddr := strings.ToLower(addr.Hex())
		if strAddr != v.addr {
			t.Errorf("TestContextApiContract creator account %v ;want %v", strAddr, v.addr)
		}
		code := stateDb.GetCode(addr)
		if wantedCode != string(code) {
			t.Errorf("TestContextApiContract code [%v] ;want [%v]", string(code), wantedCode)
		}
		iid := stateDb.GetInterpreterID(addr).Hex()
		if v.interpreter != iid {
			t.Errorf("TestContextApiContract code [%v] ;want [%v]", iid, v.interpreter)
		}
		creator := stateDb.GetCreator(addr).Hex()
		creator = strings.ToLower(creator)
		if v.creator != creator {
			t.Errorf("TestContextApiContract code [%v] ;want [%v]", creator, v.creator)
		}
	}
}
