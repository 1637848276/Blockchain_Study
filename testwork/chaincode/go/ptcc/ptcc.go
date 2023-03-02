package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// PointsTransferChaincode => 声明积分转移智能合约结构体
type PointsTransferChaincode struct {
}

// Init => 链码初始化接口
func (t *PointsTransferChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("初始化参数个数不匹配")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("积分管理员初始化失败")
	}

	return shim.Success(nil)
}

// Invoke => 链码调用接口
func (t *PointsTransferChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("无效交易方法，仅支持：transfer|query")
}

func (t *PointsTransferChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	outUser := args[0]
	inUser := args[1]
	offsetVal, _ := strconv.Atoi(args[2])

	outUserPointBalanceByte, _ := stub.GetState(outUser)
	outUserPointBalance, _ := strconv.Atoi(string(outUserPointBalanceByte))

	inUserPointBalanceByte, _ := stub.GetState(inUser)
	inUserPointBalance, _ := strconv.Atoi(string(inUserPointBalanceByte))

	outUserPointBalance = outUserPointBalance - offsetVal
	inUserPointBalance = inUserPointBalance + offsetVal

	err := stub.PutState(outUser, []byte(strconv.Itoa(outUserPointBalance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(inUser, []byte(strconv.Itoa(inUserPointBalance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *PointsTransferChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userPointsBalance, _ := stub.GetState(args[0])
	if userPointsBalance == nil {
		return shim.Error("查询账户余额失败")
	}

	return shim.Success(userPointsBalance)
}

func main() {
	err := shim.Start(new(PointsTransferChaincode))
	if err != nil {
		fmt.Printf("积分转移链码启动失败: %s\n", err)
	}
}