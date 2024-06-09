"use client";
import { useState } from "react";

const PromptInput = ({ prompt, handleInput, index, processInput, inputRefs }: { prompt: string, handleInput: (event: React.ChangeEvent<HTMLInputElement>, index: number) => void, index: number, processInput: (index: number) => void, inputRefs: React.RefObject<HTMLInputElement[]> }) => {
    const [readOnly, setReadOnly] = useState(false);
  
    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter') {
        processInput(index);
        setReadOnly(true);
  
        const nextIndex = index + 1;
        if (inputRefs.current && inputRefs.current[nextIndex]) {
          inputRefs.current[nextIndex].focus();
        }
      }
    };
  
    return (
      <div className="flex gap-2 items-center justify-start font-terminal text-black my-2">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="ml-1 w-6 h-6"
          viewBox="0 0 24 24"
        >
          <path
            fill="#5addf4"
            stroke="#22202c"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M9 10h.01M15 10h.01M12 2a8 8 0 0 0-8 8v12l3-3l2.5 2.5L12 19l2.5 2.5L17 19l3 3V10a8 8 0 0 0-8-8"
          />
        </svg>
        <h1 className="text-sky font-bold">$</h1>
        <input
          ref={el => { 
            if (inputRefs.current) {
              inputRefs.current[index] = el!;
            }
          }}
          className="border-0 outline-0 text-sky w-full mr-2 bg-terminal"
          type="text"
          value={prompt}
          readOnly={readOnly}
          onChange={(event) => handleInput(event, index)}
          onKeyDown={handleKeyDown}
        />
      </div>
    );
  };

export { PromptInput };