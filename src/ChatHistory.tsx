import React from 'react';
import { Typography, Card, CardContent } from '@mui/material';
import { connect } from './index';


interface Message {
  text: string;
  lang: string;
}

interface State {
  messages: Message[];
}

class ChatHistory extends React.Component<{}, State> {
  state: State = {
    messages: [],
  };

  componentDidMount() {
    connect((msg) => {
      console.log(msg); // Log raw message
      const data = JSON.parse(msg.data);
  
      // split the body into the translated text and original language
      const splitText = data.body.split(" [translated from ");
      const text = splitText[0];
      const lang = splitText[1] ? splitText[1].replace("]", "") : "unknown";
    
      const parsedMessage = { text, lang } as Message;
      console.log(parsedMessage); // Log parsed message
    
      this.setState((prevState) => ({
        messages: [...prevState.messages, parsedMessage],
      }));
    });
  }

  render() {
    return (
      <div>
        {this.state.messages.map((message, index) => (
          <Card key={index}>
            <CardContent>
              <Typography variant="h5">{message.text}</Typography>
              <Typography color="textSecondary">
                Translated from: {message.lang}
              </Typography>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }
}

export default ChatHistory;