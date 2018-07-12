package signer

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
This is used to generated signed messages that can be submitted to a smart contract in order to process a payment.
This module is only used for the "frontend" web GUI. Regular API use will incorporate a payment channel style contract
*/

type PaymentSigner struct {
	Key *ecdsa.PrivateKey
}

type SignedMessage struct {
	H []byte `json:"h"`
	R []byte `json:"r"`
	S []byte `json:"s"`
	V uint8  `json:"v"`
}

// GeneratePaymentSigner is used to generate our helper struct for signing payments
// keyFilePath is the path to a key as generated by geth
func GeneratePaymentSigner(keyFilePath, keyPass string) (*PaymentSigner, error) {
	fileBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	pk, err := keystore.DecryptKey(fileBytes, keyPass)
	if err != nil {
		return nil, err
	}
	return &PaymentSigner{Key: pk.PrivateKey}, nil
}

func (ps *PaymentSigner) GenerateSignedPaymentMessage(ethAddress common.Address, paymentMethod uint8, paymentNumber, chargeAmountInWei *big.Int) (*SignedMessage, error) {
	//  return keccak256(abi.encodePacked(msg.sender, _paymentNumber, _paymentMethod, _chargeAmountInWei));
	hashToSign := SoliditySHA3(
		Address(ethAddress),
		Uint256(paymentNumber),
		Uint8(paymentMethod),
		Uint256(chargeAmountInWei),
	)
	sig, err := crypto.Sign(hashToSign, ps.Key)
	if err != nil {
		return nil, err
	}
	msg := &SignedMessage{
		H: hashToSign,
		R: sig[0:32],
		S: sig[32:64],
		V: uint8(sig[64]) + 27,
	}
	return msg, nil
}
