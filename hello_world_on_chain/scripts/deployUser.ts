import { ethers } from 'hardhat';

const fs = require("fs")
const path = require("path")

async function main() {
    const contract = await ethers.deployContract('User');

    await contract.waitForDeployment();

    console.log('User Contract Deployed at ' + contract.target);

    fs.writeFile(".user_deployed_address", contract.target, err => {
        if(err) {
            console.log(err)
        }
    })
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
