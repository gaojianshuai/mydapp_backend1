// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private count;
    address public owner;
    
    event CountIncreased(uint256 newCount);
    event CountReset(uint256 newCount);

    constructor(uint256 _initialCount) {
        count = _initialCount;
        owner = msg.sender;
    }

    function increment() public {
        count += 1;
        emit CountIncreased(count);
    }

    function getCount() public view returns (uint256) {
        return count;
    }

    function reset() public {
        require(msg.sender == owner, "Only owner can reset");
        count = 0;
        emit CountReset(count);
    }
}