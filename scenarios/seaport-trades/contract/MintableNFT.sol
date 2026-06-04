// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Minimal ERC721 with permissionless, caller-chosen-id mint, used as the test
// collection for the spamoor seaport-trades scenario. Permissionless mint lets
// the scenario seed market/trader inventories without an owner role; minting a
// caller-chosen id (rather than auto-increment) lets the Go side track ids
// deterministically without parsing logs.
//
// It implements only the ERC721 surface Seaport touches when fulfilling orders
// with the zero conduit key: ownerOf, balanceOf, transferFrom, getApproved,
// setApprovalForAll/isApprovedForAll, plus approve and ERC165 supportsInterface.
// safeTransferFrom is intentionally omitted - Seaport's ERC721 transfer path
// calls transferFrom directly.
contract MintableNFT {
    string public name;
    string public symbol;

    mapping(uint256 => address) private _owners;
    mapping(address => uint256) private _balances;
    mapping(uint256 => address) private _tokenApprovals;
    mapping(address => mapping(address => bool)) private _operatorApprovals;

    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);

    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        // ERC165 (0x01ffc9a7) and ERC721 (0x80ac58cd).
        return interfaceId == 0x01ffc9a7 || interfaceId == 0x80ac58cd;
    }

    function balanceOf(address owner) external view returns (uint256) {
        require(owner != address(0), "zero address");
        return _balances[owner];
    }

    function ownerOf(uint256 tokenId) public view returns (address) {
        address owner = _owners[tokenId];
        require(owner != address(0), "nonexistent token");
        return owner;
    }

    // mint creates tokenId and assigns it to `to`. Permissionless and id-explicit
    // so the scenario controls the id space; reverts if the id already exists.
    function mint(address to, uint256 tokenId) public {
        require(to != address(0), "zero address");
        require(_owners[tokenId] == address(0), "already minted");
        _balances[to] += 1;
        _owners[tokenId] = to;
        emit Transfer(address(0), to, tokenId);
    }

    // mintBatch mints a contiguous id range [startId, startId+count) to `to`.
    function mintBatch(address to, uint256 startId, uint256 count) external {
        for (uint256 i = 0; i < count; i++) {
            mint(to, startId + i);
        }
    }

    function approve(address to, uint256 tokenId) external {
        address owner = ownerOf(tokenId);
        require(
            msg.sender == owner || _operatorApprovals[owner][msg.sender],
            "not authorized"
        );
        _tokenApprovals[tokenId] = to;
        emit Approval(owner, to, tokenId);
    }

    function getApproved(uint256 tokenId) external view returns (address) {
        require(_owners[tokenId] != address(0), "nonexistent token");
        return _tokenApprovals[tokenId];
    }

    function setApprovalForAll(address operator, bool approved) external {
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function isApprovedForAll(address owner, address operator) external view returns (bool) {
        return _operatorApprovals[owner][operator];
    }

    function transferFrom(address from, address to, uint256 tokenId) external {
        address owner = ownerOf(tokenId);
        require(owner == from, "wrong from");
        require(to != address(0), "zero address");
        require(
            msg.sender == owner ||
                _tokenApprovals[tokenId] == msg.sender ||
                _operatorApprovals[owner][msg.sender],
            "not authorized"
        );

        delete _tokenApprovals[tokenId];
        _balances[from] -= 1;
        _balances[to] += 1;
        _owners[tokenId] = to;
        emit Transfer(from, to, tokenId);
    }
}
