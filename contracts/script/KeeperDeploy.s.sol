// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Script, console} from "forge-std/Script.sol";
import {Keeper} from "../src/Keeper.sol";

contract KeeperDeployScript is Script {
    function run() public {
        vm.startBroadcast();
        Keeper keeper = new Keeper();
        console.log("Keeper deployed at:", address(keeper));
        vm.stopBroadcast();
    }
}
