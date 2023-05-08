package ensoul

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/0xDeSchool/gap/app"
	"math/big"

	abi "github.com/0xDeSchool/gap/ensoul/abi"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Sdk struct {
	Opts      *EnsoulSdkOptions
	EthClient *ethclient.Client
}

type EnsoulSdkOptions struct {
	PoaPrivKey           string
	PoskPrivKey          string
	DescPoaContractAddr  string
	DescPoskContractAddr string
	ChainId              int
}

func NewSdk() *Sdk {
	opts := app.Get[EnsoulSdkOptions]()
	ethc := app.Get[ethclient.Client]()
	ensoulSdk := &Sdk{
		Opts:      opts,
		EthClient: ethc,
	}
	return ensoulSdk
}

func (es *Sdk) MintOnePoAToken(ctx context.Context, toAddrStr string, tokenIdStr string) *types.Transaction {
	poaAddr := common.HexToAddress(es.Opts.DescPoaContractAddr)
	return es.MintOneToken(ctx, poaAddr, toAddrStr, tokenIdStr, "poa")
}

func (es *Sdk) MintOnePoSkToken(ctx context.Context, toAddrStr string, tokenIdStr string) *types.Transaction {
	poskAddr := common.HexToAddress(es.Opts.DescPoskContractAddr)
	return es.MintOneToken(ctx, poskAddr, toAddrStr, tokenIdStr, "posk")
}

func (es *Sdk) MintOneToken(ctx context.Context, contractAddr common.Address, toAddrStr string, tokenIdStr string, contractType string) *types.Transaction {
	toAddr := common.HexToAddress(toAddrStr)
	tx := es.SendEnsoulToken(ctx, contractAddr, toAddr, tokenIdStr, contractType)
	// 目前设定并发监听一下
	go es.SubscribePendingTx(ctx, tx)
	return tx
}

func (es *Sdk) SendEnsoulToken(ctx context.Context, contractAddr common.Address, toAddr common.Address, tokenIdStr string, contractType string) *types.Transaction {
	// 准备 Ensoul 合约
	ensoulV11, err := abi.NewEnsoul(contractAddr, es.EthClient)
	errx.CheckError(err)

	// 获取 Auth
	var secret string
	if contractType == "poa" {
		secret = es.Opts.PoaPrivKey
	} else {
		secret = es.Opts.PoskPrivKey
	}
	auth := es.PrepareTxAuth(secret)

	// 设置数据
	tokenIdBig := new(big.Int)
	tokenIdBig, ok := tokenIdBig.SetString(tokenIdStr, 10)
	if !ok {
		errx.Panic("Invalid tokenId given")
	}
	fmt.Println(tokenIdBig)
	// 一般都是固定数额的 SBT，一枚
	amountBig := big.NewInt(1)

	// 发送与监听交易
	tx, err := ensoulV11.Mint(auth, toAddr, tokenIdBig, amountBig)
	errx.CheckError(err)
	log.Infof("Tx to ensoul sent: %s\n", tx.Hash().Hex())
	return tx
}

func (es *Sdk) SubscribePendingTx(ctx context.Context, tx *types.Transaction) bool {
	// 输出 tx 并等待挖出
	bind.WaitMined(context.Background(), es.EthClient, tx)
	log.Infof("Tx to ensoul minted: %s\n", tx.Hash().Hex())

	return true
}

func (es *Sdk) PrepareTxAuth(privKeyStr string) *bind.TransactOpts {
	// 准备 Signer: PublicKey, SignerAddress
	privKey, err := crypto.HexToECDSA(privKeyStr)
	errx.CheckError(err)
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		errx.Panic("unable to generate PubKey")
	}
	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 准备 网络信息
	nonce, err := es.EthClient.PendingNonceAt(context.Background(), signerAddress)
	errx.CheckError(err)

	gasPrice, err := es.EthClient.SuggestGasPrice(context.Background())
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))
	errx.CheckError(err)

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(int64(es.Opts.ChainId)))
	errx.CheckError(err)

	// 设置参数
	auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(1000000000) // in wei
	// auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}
