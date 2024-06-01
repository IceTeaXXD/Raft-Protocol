'use client';

import { useState, useRef, useEffect } from 'react';

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
    <div className="flex gap-2 items-center justify-start font-ProFontIIxNerdFontRegular text-white my-2">
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
        ref={el => { 
          if (inputRefs.current) {
            inputRefs.current[index] = el!;
          }
        }}
        className="border-0 outline-0 text-sm text-black w-full mr-2"
        type="text"
        value={prompt}
        readOnly={readOnly}
        onChange={(event) => handleInput(event, index)}
        onKeyDown={handleKeyDown}
      />
    </div>
  );
};

const Home = () => {
  const [prompts, setPrompts] = useState([{ input: '', output: ''}]);
  const inputRefs = useRef<HTMLInputElement[]>([]);

  const handleInput = (event: React.ChangeEvent<HTMLInputElement>, index: number) => {
    const newPrompts = [...prompts];
    newPrompts[index].input = event.target.value;
    setPrompts(newPrompts);
  };

  const processInput = (index: number) => {
    const processedInput = prompts[index].input;
    const newPrompts = [...prompts];
    newPrompts[index].output = processedInput;
    setPrompts([...newPrompts, { input: '', output: '' }]);
  };

  useEffect(() => {
    if (inputRefs.current.length > 0) {
      inputRefs.current[inputRefs.current.length - 1].focus();
    }
  }, [prompts]);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="w-[96%] h-[90%] md:w-[700px] md:h-[500px] lg:w-[800px] lg:h-[600px] flex flex-col border-8 border-black rounded-lg bg-black">
          <div className="relative flex flex-col items-center">
          </div>
          <div className="w-full flex-1 bg-white pt-2 overflow-y-scroll scrollbar">
            <h1 className="hidden md:block pl-2 my-2 font-profontiixnerdfont text-black text-sm">
              Test
            </h1>
            <div>
              {prompts.map((prompt, index) => (
                <div key={index}>
                  <PromptInput prompt={prompt.input} handleInput={handleInput} index={index} processInput={processInput} inputRefs={inputRefs} />
                  {prompt.output && (
                    <div className="pl-2 text-black font-profontiixnerdfont w-full break-words">
                      Processed command: {prompt.output}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
    </main>
  );
};

export default Home;
