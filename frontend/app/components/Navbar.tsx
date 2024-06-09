import React from 'react'

function Navbar() {
  return (
    <nav className="bg-terminal p-4">
    <ul className="flex space-x-4">
      <li>
        <a href="/" className="text-sky font-terminal font-bold hover:text-gray-400">SPG</a>
      </li>
      <li>
        <a href="/dashboard" className="text-sky font-terminal hover:text-gray-400">Dashboard</a>   
      </li>
      <li>
        <a href="/terminal" className="text-sky font-terminal hover:text-gray-400">Terminal</a>
      </li>
    </ul>
  </nav>
  )
}

export default Navbar