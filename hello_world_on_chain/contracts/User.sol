// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParser} from "./TAParser.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable, TAParser {
    event AttestedFunctionCallOutput(string output);

    address private taSigningKeyAddress;

    constructor() Ownable(msg.sender) {
    }

    function setTASigningKeyAddress(
        bytes calldata taSigningKey
    )
    public onlyOwner
    {
        taSigningKeyAddress = publicKeyToAddress(taSigningKey);
    }

    function verifyAttestedFnCallClaims(
        string calldata taData
    )
    public
    {
        TAParser.FnCallClaims memory claims = parseTA(
            taData,
            taSigningKeyAddress
        );

        console.log("Verified attest-fn-call claims:");
        console.log("\tFunction: %s", claims.Function);
        console.log("\tHash of code: %s", claims.HashOfCode);
        console.log("\tHash of input: %s", claims.HashOfInput);
        console.log("\tHash of secrets: %s", claims.HashOfSecrets);
        console.log("\tOutput: %s", claims.Output);

        emit AttestedFunctionCallOutput(claims.Output);
    }
}

