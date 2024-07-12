require("@nomiclabs/hardhat-ethers");
require("@nomicfoundation/hardhat-chai-matchers");

require("dotenv").config({ path: "../../../.env" });
require("./tasks/proveComet");

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.8.19",
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {},
    sepolia: {
      url: "https://rpc-sepolia-eth.nodeguardians.io",
    },
  },
  paths: {
    sources: "./contracts_",
  },
};
