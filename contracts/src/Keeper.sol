// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IKeeper} from "./Interfaces/IKeeper.sol";

error InvalidDataLength();
error CannotStoreExistingData(address, uint256);
error CannotChangeNonExistentData(address, uint256);
error CannotRemoveNonExistentData(address, uint256);

contract Keeper is IKeeper {
    mapping(address => mapping(uint256 => bytes)) public userData;
    mapping(address => bytes) public userMetaData;
    mapping(address => uint256[]) public activeIdsForUser;
    mapping(address => uint256) public nextDataId;

    function storeMetaData(bytes calldata _data) external payable {
        require(_data.length > 0, InvalidDataLength());
        address account = msg.sender;
        userMetaData[account] = _data;
    }

    function storeData(bytes calldata _data) external payable {
        require(_data.length > 0, InvalidDataLength());
        address account = msg.sender;
        uint256 id = nextDataId[account];
        require(userData[account][id].length == 0, CannotStoreExistingData(account, id));
        nextDataId[account]++;
        activeIdsForUser[account].push(id);

        userData[account][id] = _data;
    }

    function changeData(uint256 _id, bytes calldata _newData) external payable {
        require(_newData.length > 0, InvalidDataLength());
        address account = msg.sender;
        require(userData[account][_id].length != 0, CannotChangeNonExistentData(account, _id));

        userData[account][_id] = _newData;
    }

    function removeData(uint256 _id) external payable {
        address account = msg.sender;
        require(userData[account][_id].length != 0, CannotRemoveNonExistentData(account, _id));

        delete userData[account][_id];

        uint256[] storage activeIds = activeIdsForUser[account];
        uint256 length = activeIds.length;
        for (uint256 i = 0; i < length;) {
            if (activeIds[i] == _id) {
                activeIds[i] = activeIds[length - 1];
                activeIds.pop();
                break;
            }
            unchecked {
                ++i;
            }
        }
    }
}
