package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// SEPA XML structures for pain.001.001.03 format
type Document struct {
	XMLName                xml.Name               `xml:"Document"`
	CustomerCreditTransfer CustomerCreditTransfer `xml:"CstmrCdtTrfInitn"`
}

type CustomerCreditTransfer struct {
	GroupHeader GroupHeader `xml:"GrpHdr"`
	PaymentInfo PaymentInfo `xml:"PmtInf"`
}

type GroupHeader struct {
	MessageID        string    `xml:"MsgId"`
	CreationDateTime string    `xml:"CreDtTm"`
	NumberOfTxs      string    `xml:"NbOfTxs"`
	ControlSum       string    `xml:"CtrlSum"`
	InitiatingParty  PartyInfo `xml:"InitgPty"`
}

type PaymentInfo struct {
	PaymentInfoID        string               `xml:"PmtInfId"`
	PaymentMethod        string               `xml:"PmtMtd"`
	BatchBooking         string               `xml:"BtchBookg"`
	NumberOfTxs          string               `xml:"NbOfTxs"`
	ControlSum           string               `xml:"CtrlSum"`
	PaymentTypeInfo      PaymentTypeInfo      `xml:"PmtTpInf"`
	RequestedExecDate    string               `xml:"ReqdExctnDt"`
	Debtor               PartyInfo            `xml:"Dbtr"`
	DebtorAccount        Account              `xml:"DbtrAcct"`
	DebtorAgent          FinancialInstitution `xml:"DbtrAgt"`
	ChargeBearer         string               `xml:"ChrgBr"`
	CreditTransferTxInfo CreditTransferTxInfo `xml:"CdtTrfTxInf"`
}

type PaymentTypeInfo struct {
	ServiceLevel    ServiceLevel    `xml:"SvcLvl"`
	CategoryPurpose CategoryPurpose `xml:"CtgyPurp"`
}

type ServiceLevel struct {
	Code string `xml:"Cd"`
}

type CategoryPurpose struct {
	Code string `xml:"Cd"`
}

type PartyInfo struct {
	Name string `xml:"Nm"`
}

type Account struct {
	ID       AccountID `xml:"Id"`
	Currency string    `xml:"Ccy"`
}

type AccountID struct {
	IBAN string `xml:"IBAN"`
}

type FinancialInstitution struct {
	FinInstnID FinInstnID `xml:"FinInstnId"`
}

type FinInstnID struct {
	BIC string `xml:"BIC"`
}

type CreditTransferTxInfo struct {
	PaymentID       PaymentID            `xml:"PmtId"`
	Amount          Amount               `xml:"Amt"`
	CreditorAgent   FinancialInstitution `xml:"CdtrAgt"`
	Creditor        PartyInfo            `xml:"Cdtr"`
	CreditorAccount Account              `xml:"CdtrAcct"`
	RemittanceInfo  RemittanceInfo       `xml:"RmtInf"`
}

type PaymentID struct {
	EndToEndID string `xml:"EndToEndId"`
}

type Amount struct {
	InstructedAmount InstructedAmount `xml:"InstdAmt"`
}

type InstructedAmount struct {
	Currency string `xml:"Ccy,attr"`
	Value    string `xml:",chardata"`
}

type RemittanceInfo struct {
	Structured StructuredRemittance `xml:"Strd"`
}

type StructuredRemittance struct {
	CreditorRefInfo CreditorRefInfo `xml:"CdtrRefInf"`
}

type CreditorRefInfo struct {
	Type      RefType `xml:"Tp"`
	Reference string  `xml:"Ref"`
}

type RefType struct {
	CodeOrProprietary CodeOrProprietary `xml:"CdOrPrtry"`
	Issuer            string            `xml:"Issr"`
}

type CodeOrProprietary struct {
	Code string `xml:"Cd"`
}

// SEPAData represents parsed SEPA data in a flat structure for table display
type SEPAData struct {
	Fields []SEPAField
}

type SEPAField struct {
	Category string
	Field    string
	Value    string
}

// ParseSEPAFile parses a SEPA XML file and returns structured data
func ParseSEPAFile(filepath string) (*SEPAData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return convertToSEPAData(&doc), nil
}

// convertToSEPAData converts the parsed XML structure to a flat table structure
func convertToSEPAData(doc *Document) *SEPAData {
	var fields []SEPAField

	// Group Header
	fields = append(fields, SEPAField{"Group Header", "Message ID", doc.CustomerCreditTransfer.GroupHeader.MessageID})
	fields = append(fields, SEPAField{"Group Header", "Creation Date Time", formatDateTime(doc.CustomerCreditTransfer.GroupHeader.CreationDateTime)})
	fields = append(fields, SEPAField{"Group Header", "Number of Transactions", doc.CustomerCreditTransfer.GroupHeader.NumberOfTxs})
	fields = append(fields, SEPAField{"Group Header", "Control Sum", doc.CustomerCreditTransfer.GroupHeader.ControlSum})
	fields = append(fields, SEPAField{"Group Header", "Initiating Party", doc.CustomerCreditTransfer.GroupHeader.InitiatingParty.Name})

	// Payment Info
	pmtInfo := doc.CustomerCreditTransfer.PaymentInfo
	fields = append(fields, SEPAField{"Payment Info", "Payment Info ID", pmtInfo.PaymentInfoID})
	fields = append(fields, SEPAField{"Payment Info", "Payment Method", pmtInfo.PaymentMethod})
	fields = append(fields, SEPAField{"Payment Info", "Batch Booking", pmtInfo.BatchBooking})
	fields = append(fields, SEPAField{"Payment Info", "Number of Transactions", pmtInfo.NumberOfTxs})
	fields = append(fields, SEPAField{"Payment Info", "Control Sum", pmtInfo.ControlSum})
	fields = append(fields, SEPAField{"Payment Info", "Service Level", pmtInfo.PaymentTypeInfo.ServiceLevel.Code})
	fields = append(fields, SEPAField{"Payment Info", "Category Purpose", pmtInfo.PaymentTypeInfo.CategoryPurpose.Code})
	fields = append(fields, SEPAField{"Payment Info", "Requested Execution Date", pmtInfo.RequestedExecDate})
	fields = append(fields, SEPAField{"Payment Info", "Charge Bearer", pmtInfo.ChargeBearer})

	// Debtor Info
	fields = append(fields, SEPAField{"Debtor", "Name", pmtInfo.Debtor.Name})
	fields = append(fields, SEPAField{"Debtor", "IBAN", pmtInfo.DebtorAccount.ID.IBAN})
	fields = append(fields, SEPAField{"Debtor", "Currency", pmtInfo.DebtorAccount.Currency})
	fields = append(fields, SEPAField{"Debtor", "BIC", pmtInfo.DebtorAgent.FinInstnID.BIC})

	// Credit Transfer Transaction Info
	txInfo := pmtInfo.CreditTransferTxInfo
	fields = append(fields, SEPAField{"Transaction", "End to End ID", txInfo.PaymentID.EndToEndID})
	fields = append(fields, SEPAField{"Transaction", fmt.Sprintf("Amount (%s)", txInfo.Amount.InstructedAmount.Currency), txInfo.Amount.InstructedAmount.Value})

	// Creditor Info
	fields = append(fields, SEPAField{"Creditor", "Name", txInfo.Creditor.Name})
	fields = append(fields, SEPAField{"Creditor", "IBAN", txInfo.CreditorAccount.ID.IBAN})
	fields = append(fields, SEPAField{"Creditor", "BIC", txInfo.CreditorAgent.FinInstnID.BIC})

	// Remittance Info
	fields = append(fields, SEPAField{"Remittance", "Reference Type", txInfo.RemittanceInfo.Structured.CreditorRefInfo.Type.CodeOrProprietary.Code})
	fields = append(fields, SEPAField{"Remittance", "Issuer", txInfo.RemittanceInfo.Structured.CreditorRefInfo.Type.Issuer})
	fields = append(fields, SEPAField{"Remittance", "Reference", txInfo.RemittanceInfo.Structured.CreditorRefInfo.Reference})

	return &SEPAData{Fields: fields}
}

// formatDateTime formats the ISO datetime string to a more readable format
func formatDateTime(dateTime string) string {
	if t, err := time.Parse("2006-01-02T15:04:05", dateTime); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return dateTime
}
