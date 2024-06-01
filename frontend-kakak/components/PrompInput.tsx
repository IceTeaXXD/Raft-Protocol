// components/PromptInput.js
import { useState } from 'react';

const PromptInput = ({ prompt, handleInput, index }) => {
  return (
    <div className="flex gap-2 items-center justify-start font-ProFontIIxNerdFontRegular text-white">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className="ml-1 w-6 h-6"
        viewBox="0 0 24 24"
      >
        <path
          fill="#c6aae8"
          stroke="#22202c"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="1"
          d="M9 10h.01M15 10h.01M12 2a8 8 0 0 0-8 8v12l3-3l2.5 2.5L12 19l2.5 2.5L17 19l3 3V10a8 8 0 0 0-8-8"
        />
      </svg>
      <h1 className="text-black text-sm font-bold">$</h1>
      <input
        className="border-0 outline-0 text-sm text-black w-full mr-2 placeholder-gray-500"
        type="text"
        placeholder="Enter text here..."
        value={prompt}
        onChange={(event) => handleInput(event, index)}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            // Call processInput function from the parent component
            processInput(index);
          }
        }}
      />
    </div>
  );
};

export default PromptInput;
