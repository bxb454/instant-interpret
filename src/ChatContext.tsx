import React, { createContext, useState, ReactNode } from 'react';

interface Message {
  text: string;
  lang: string;
}

interface ChatContextInterface {
  messages: Message[];
  addMessage: (message: Message) => void;
}

export const ChatContext = createContext<ChatContextInterface>({
  messages: [],
  addMessage: (message: Message) => {},
});

interface ChatProviderProps {
  children: ReactNode;
}

const ChatProvider: React.FC<ChatProviderProps> = ({ children }) => {
  const [messages, setMessages] = useState<Message[]>([]);

  const addMessage = (message: Message) => {
    setMessages((prevMessages) => [...prevMessages, message]);
  };

  return (
    <ChatContext.Provider value={{ messages, addMessage }}>
      {children}
    </ChatContext.Provider>
  );
};

export default ChatProvider;