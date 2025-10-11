// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface IKeeper {
    function storeMetaData(bytes calldata _data) external payable;
    function storeData(bytes calldata _data) external payable;
    function changeData(uint256 _id, bytes calldata _newData) external payable;
    function removeData(uint256 _id) external payable;
}
