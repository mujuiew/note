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
		return
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("2Failed to read file: %v", err)
		return
	}

	// Parse the JSON array
	var body Body
	err = json.Unmarshal(content, &body)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Add header to output
	out = append(out, "AccountNumber|EventCode|PrincipalAmount|InterestAmount|PenaltyAmount|VatAmount|FeeAmount|otherProperties[bill]|otherProperties[penalty]|otherProperties[advance_payment]|otherProperties[fee]")

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

			// Combine bills and penalties for otherProperties
			var bill string
			var penalty string
			if rawBills, ok := dueBillsMsg.OtherProperties["bills"]; ok {
				if billsStr, ok := rawBills.(string); ok {
					bill = billsStr
				}
			}

			// Check and assign "penalty" field
			if rawPenalty, ok := dueBillsMsg.OtherProperties["penalty"]; ok {
				if penaltyStr, ok := rawPenalty.(string); ok {
					penalty = penaltyStr
				}
			}

			out = append(out, fmt.Sprintf("%s|%s|%s|%s|%s|%s||%s|%s||",
				tx.AccountNumber.String(), tx.EventCode.String(), dueBillsMsg.PrincipalAmount.String(),
				dueBillsMsg.InterestAmount.String(), dueBillsMsg.PenaltyAmount.String(),
				dueBillsMsg.VatAmount.String(), bill, penalty))

		case "fee":
			var feeMsg FeeMessage
			err := json.Unmarshal([]byte(tx.Message.String()), &feeMsg)
			if err != nil {
				fmt.Println("Error unmarshalling fee message:", err)
				continue
			}
			// Include fee details in otherProperties

			var fee string
			if rawBills, ok := feeMsg.OtherProperties["fee"]; ok {
				if billsStr, ok := rawBills.(string); ok {
					fee = billsStr
				}
			}
			out = append(out, fmt.Sprintf("%s|%s|||||%s||||%s",
				tx.AccountNumber.String(), tx.EventCode.String(), feeMsg.FeeAmount.String(), fee))

		case "others":
			var othersMsg OthersMessage
			err := json.Unmarshal([]byte(tx.Message.String()), &othersMsg)
			if err != nil {
				fmt.Println("Error unmarshalling others message:", err)
				continue
			}
			// Combine advance payment and penalties for otherProperties

			var advancePayment string
			var penalty string
			if rawBills, ok := othersMsg.OtherProperties["advance_payment"]; ok {
				if billsStr, ok := rawBills.(string); ok {
					advancePayment = billsStr
				}
			}

			// Check and assign "penalty" field
			if rawPenalty, ok := othersMsg.OtherProperties["penalty"]; ok {
				if penaltyStr, ok := rawPenalty.(string); ok {
					penalty = penaltyStr
				}
			}
			out = append(out, fmt.Sprintf("%s|%s|%s|%s|%s|%s|||%s|%s|",
				tx.AccountNumber.String(), tx.EventCode.String(), othersMsg.PrincipalAmount.String(), othersMsg.InterestAmount.String(), othersMsg.PenaltyAmount.String(),
				othersMsg.VatAmount.String(), penalty, advancePayment))

		default:
			fmt.Printf("Unknown EventCode: %s\n", tx.EventCode)
		}
	}

	// Print the output
	for _, line := range out {
		fmt.Println(line)
	}
}
