import { TextField, Button, Select, MenuItem } from '@mui/material';
import { ChatContext } from './ChatContext';
import { sendMsg } from "./index";
import { useEffect, useContext, useState } from 'react';

const languages = ["en", "es", "fr", "de"];
const MsgTypeUserMessage = 1;
const MsgTypeLangUpdate = 2;

const ChatMsg = () => {
  const { addMessage } = useContext(ChatContext);
  const [message, setMessage] = useState("");
  const [language, setLanguage] = useState(languages[0]);

  useEffect(() => {
    // When language changes, send a language update message
    const messageObject = { type: MsgTypeLangUpdate, body: language };
    sendMsg(JSON.stringify(messageObject));
  }, [language]);

  const handleSend = () => {
    // When sending a message, send it with type user message and the current language
    const messageObject = { type: MsgTypeUserMessage, body: message, originalLanguage: language };
    sendMsg(JSON.stringify(messageObject));
    addMessage({text: message, lang: language});
    setMessage("");
  };

  const handleDelete = () => {
    setMessage("");
  };

  return (
    <div>
      <TextField 
        label="Message"
        variant="outlined"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <Select
        value={language}
        onChange={(e) => setLanguage(e.target.value)}
      >
        {languages.map((lang, index) => (
          <MenuItem key={index} value={lang}>{lang}</MenuItem>
        ))}
      </Select>
      <Button variant="contained" color="primary" onClick={handleSend}>
        Send
      </Button>
      <Button variant="contained" color="secondary" onClick={handleDelete}>
        Delete
      </Button>
    </div>
  );
};

export default ChatMsg;