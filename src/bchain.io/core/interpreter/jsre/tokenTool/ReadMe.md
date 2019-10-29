The purpose of this tool is modifying the templateJs.go's content.
The result (output file),is the token contract code.

If you want publish a Token on the BchainChain.Follow steps bellow.

1.config
name: test.config

	{
	"symbol":"ABC",
	"token_name":"ABC Token",
	"decimals":"18",
	"total_supply":"1e+10",
	"times_total_supply":"1e+18",
	"bchain_contract":"0x192d52D8cE0c7bBAf0780EAb04860D6Ba012578B"
	}

2.use this tool binary file to modify the templateJs

tokenTool  test.config  output.md

3.copy the content of output.md , sign a transaction to publish a token contract.

*******************************
Notice:
your current path should contain the templateJsCode.js file
Because if you want modify or add content of the contract , you can modify templateJsCode directly,
and after that,you use the step upon to generate the output.md,the result file you want.