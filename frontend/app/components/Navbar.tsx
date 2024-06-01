import React from 'react'

function Navbar() {
  return (
    <nav className="bg-black p-4">
    <ul className="flex space-x-4">
      <li>
        <a href="/" className="text-white font-bold hover:text-gray-400">SPG</a>
      </li>
      <li>
        <a href="/dashboard" className="text-white hover:text-gray-400">Dashboard</a>   
      </li>
      <li>
        <a href="/terminal" className="text-white hover:text-gray-400">Terminal</a>
      </li>
    </ul>
  </nav>
  )
}

export default Navbar