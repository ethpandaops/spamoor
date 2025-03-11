// SPDX-License-Identifier: MIT
pragma solidity  ^0.8.22;

interface IDepositContract {
    function deposit(bytes calldata pubkey, bytes calldata withdrawal_credentials, bytes calldata signature, bytes32 deposit_data_root) external payable;
}

contract BatchTopUps {
    address public _depositContract;

    constructor(address depositContract) payable {
        _depositContract = depositContract;
    }

    function _topup(bytes calldata pubkey, uint256 amount) internal {
        bytes32 withdrawal_credentials = 0x0000000000000000000000000000000000000000000000000000000000000000;
        bytes memory amount_bytes = to_little_endian_64(uint64(amount / 1 gwei));
        bytes memory signature = abi.encodePacked(bytes32(0), bytes32(0), bytes32(0));
        bytes32 pubkey_root = sha256(abi.encodePacked(pubkey, bytes16(0)));
        bytes32 signature_root = sha256(abi.encodePacked(
            sha256(abi.encodePacked(bytes32(0), bytes32(0))),
            sha256(abi.encodePacked(bytes32(0), bytes32(0)))
        ));
        bytes32 deposit_data_root = sha256(abi.encodePacked(
            sha256(abi.encodePacked(pubkey_root, withdrawal_credentials)),
            sha256(abi.encodePacked(amount_bytes, bytes24(0), signature_root))
        ));
        IDepositContract(_depositContract).deposit{value: amount}(pubkey, abi.encodePacked(withdrawal_credentials), signature, deposit_data_root);
    }

    function to_little_endian_64(uint64 value) internal pure returns (bytes memory ret) {
        ret = new bytes(8);
        bytes8 bytesValue = bytes8(value);
        // Byteswapping during copying to bytes.
        ret[0] = bytesValue[7];
        ret[1] = bytesValue[6];
        ret[2] = bytesValue[5];
        ret[3] = bytesValue[4];
        ret[4] = bytesValue[3];
        ret[5] = bytesValue[2];
        ret[6] = bytesValue[1];
        ret[7] = bytesValue[0];
    }

    function topupEqual(bytes calldata pubkeys) public payable {
        uint32 pubkeysLen = uint32(pubkeys.length);
        uint256 amount = msg.value * 48 / pubkeysLen;
        uint32 pos = 0;
        require(amount >= 1 ether, "amount too low");

        while (pos < pubkeysLen) {
            _topup(pubkeys[pos:pos+48], amount);
            unchecked { pos += 48; }
        }
    }

    function topup(bytes calldata pubkey) public payable {
        uint256 amount = msg.value;
        require(amount >= 1 ether, "amount too low");
        _topup(pubkey, amount);
    }

}
