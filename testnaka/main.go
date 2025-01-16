package testnaka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Main2() {
	var out []string
	file, err := os.Open("query-dloan-payment-publishMessageDetail_response.json")
	if err != nil {
		fmt.Printf("1Failed to read file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("2Failed to read file: %v", err)
	}

	// Parse the JSON array
	var body Body
	err = json.Unmarshal(content, &body)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	out = append(out, "AccountNumber|EventCode|PrincipalAmount|InterestAmount|PenaltyAmount|VatAmount|FeeAmount\n")
	// Process each transaction
	for _, tx := range body.ReqBody {

		switch tx.EventCode.String() {
		case "due_bills":
			var dueBillsMsg DueBillsMessage
			err := json.Unmarshal([]byte(tx.Message.String()), &dueBillsMsg)
			if err != nil {
				fmt.Println("Error unmarshalling due_bills message:", err)
				continue
			}
			out = append(out, tx.AccountNumber.String()+"|"+tx.EventCode.Val+"|"+dueBillsMsg.PrincipalAmount.Val.String()+"|"+dueBillsMsg.InterestAmount.Val.String()+"|"+dueBillsMsg.PenaltyAmount.Val.String()+"|"+dueBillsMsg.VatAmount.Val.String()+"|\n")

		case "fee":
			var feeMsg FeeMessage
			err := json.Unmarshal([]byte(tx.Message.String()), &feeMsg)
			if err != nil {
				fmt.Println("Error unmarshalling fee message:", err)
				continue
			}
			out = append(out, tx.AccountNumber.String()+"|"+tx.EventCode.Val+"|||||"+feeMsg.FeeAmount.Val.String()+"\n")

		case "others":
			var othersMsg OthersMessage
			err := json.Unmarshal([]byte(tx.Message.String()), &othersMsg)
			if err != nil {
				fmt.Println("Error unmarshalling others message:", err)
				continue
			}
			out = append(out, tx.AccountNumber.String()+"|"+tx.EventCode.Val+"|"+othersMsg.PrincipalAmount.Val.String()+"|"+othersMsg.InterestAmount.Val.String()+"|"+othersMsg.PenaltyAmount.Val.String()+"|"+othersMsg.VatAmount.Val.String()+"|\n")

		default:
			fmt.Printf("Unknown EventCode: %s\n", tx.EventCode)
		}
	}

	fmt.Println(out)
}
