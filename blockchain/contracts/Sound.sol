pragma solidity ^0.5.7;

contract Sound{
    int[] sonido;
    uint precio;
    int rate;
    address[] compradores;
    address payable private owner;

    constructor(uint _precio,int _rate) public {
        owner = msg.sender;
        precio = _precio;
        rate = _rate;
    }

    function upload (int[] memory audioPart) public{
        if(msg.sender != owner)revert("No está autorizado");
        for(uint i = 0; i < audioPart.length;i++)
            sonido.push(audioPart[i]);
    }
    
    function find(address[] memory v,address e) private pure returns(bool)
    {
        bool found = false;

        for(uint i = 0; i < v.length && !found; i++)
            found = v[i] == e;

        return found;
    }

    function getLen() public view returns(uint)
    {
        return sonido.length;
    }

    function download (uint first, uint last) public view returns (int[] memory)
    {
        if(!find(compradores,msg.sender)) revert("No ha comprado el audio aun");
        int[] memory aux = new int[](last-first+1);

        for(uint i = first; i <= last; i++)
        {
            aux[i-first] = sonido[i];
        }
        
        return aux;
    }

    function getRate () public view returns(int)
    {
        return rate;
    }

    function buySound() public payable{
        if(msg.value < precio )revert("No es suficiente dinero");
        if(!find(compradores,msg.sender))
        {
            owner.transfer(precio);
            compradores.push(msg.sender);
            sonido.push(rate);
        }
    }

    function compras() public view returns(address[] memory)
    {
        return compradores;
    }

    function getPrecio() public view returns(uint)
    {
        return precio;
    }
}