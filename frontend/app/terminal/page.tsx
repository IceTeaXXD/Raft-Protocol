"use client";

import Navbar from "../components/Navbar";
import { PromptInput } from "../components/PromptInput";
import { useState, useRef, useEffect } from "react";
import axios from "axios";

function TerminalPage() {
  const [prompts, setPrompts] = useState([{ input: "", output: "" }]);
  const inputRefs = useRef<HTMLInputElement[]>([]);

  const handleInput = (
    event: React.ChangeEvent<HTMLInputElement>,
    index: number
  ) => {
    const newPrompts = [...prompts];
    newPrompts[index].input = event.target.value;
    setPrompts(newPrompts);
  };

  const processInput = async (index: number) => {
    const processedInput = prompts[index].input.trim();
    const [command, key, value] = processedInput.split(" ");
    let output = "";
    const host = process.env.NEXT_PUBLIC_SERVER_HOST || "localhost";
    const port = process.env.NEXT_PUBLIC_SERVER_PORT || "8080";
    try {
      if (command === "ping") {
        const response = await axios.get(`http://${host}:${port}/ping`);
        output = response.data.response;
      } else if (
        command === "set" &&
        key !== undefined &&
        value !== undefined
      ) {
        const response = await axios.put(
          `http://${host}:${port}/set?key=${key}&value=${value}`
        );
        output = response.data.response;
        if (output === "") {
          output = `""`;
        }
      } else if (
        command === "append" &&
        key !== undefined &&
        value !== undefined
      ) {
        const response = await axios.put(
          `http://${host}:${port}/append?key=${key}&value=${value}`
        );
        output = response.data.response;
        if (output === "") {
          output = `""`;
        }
      } else if (command === "get" && key !== undefined) {
        const response = await axios.get(
          `http://${host}:${port}/get?key=${key}`
        );
        output = response.data.response;
        console.log(output);
        output = output.toString();
        console.log(output);
        if (output === "") {
          output = `""`;
        }
      } else if (command === "strln" && key !== undefined) {
        const response = await axios.get(
          `http://${host}:${port}/strln?key=${key}`
        );
        let len = response.data.response;
        output = len.toString();
        console.log(output);
        if (output === "") {
          output = `""`;
        }
      } else if (command === "del" && key !== undefined) {
        const response = await axios.delete(
          `http://${host}:${port}/del?key=${key}`
        );
        output = response.data.response;
        if (output === "") {
          output = `""`;
        }
      } else if (command === "help") {
        output = `Usage:\n 1. ping\n 2. set <key> <value>\n 3. append <key> <value>\n 4. get <key>\n 5. strln <key>\n 6. del <key>`;
      } else {
        output = 'Invalid command, type "help" for usage';
      }
    } catch (error) {
      output = `Error: ${error}`;
    }

    const newPrompts = [...prompts];
    newPrompts[index].output = output;
    setPrompts([...newPrompts, { input: "", output: "" }]);
  };

  useEffect(() => {
    if (inputRefs.current.length > 0) {
      inputRefs.current[inputRefs.current.length - 1].focus();
    }
  }, [prompts]);

  return (
    <div>
      <Navbar />
      <main className="flex flex-col items-center justify-between p-24">
        <div className="w-[96%] h-[90%] md:w-[700px] md:h-[500px] lg:w-[800px] lg:h-[600px] flex flex-col border-8 border-terminal rounded-lg bg-terminal">
          <div className="w-full flex-1 bg-terminal pt-1 overflow-y-scroll scrollbar">
            <h1 className="hidden md:block pr-2 font-terminal text-sky text-right">
              {">"}_
            </h1>
            <div>
              {prompts.map((prompt, index) => (
                <div key={index}>
                  <PromptInput
                    prompt={prompt.input}
                    handleInput={handleInput}
                    index={index}
                    processInput={processInput}
                    inputRefs={inputRefs}
                  />
                  {prompt.output && (
                    <div className="pl-2 text-sky font-terminal w-full break-words">
                      {prompt.output.split("\n").map((line, index) => (
                        <div key={index}>{line}</div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default TerminalPage;
