import React, { Component } from "react";
import { hot } from "react-hot-loader";
import "./App.css";

class App extends Component {
    constructor() {
        super();
        this.state = {
            result : '',
            error: ''
        };
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(event) {
        event.preventDefault();
        const data = new FormData(event.target);

        if (data.get('number').length == 0) {
            this.setState( {error : 'Phone number cannot be empty'});
            return;
        }

        if (data.get('text1').length == 0 && data.get('text2').length == 0 && data.get('text3').length == 0) {
            this.setState( {error : 'Please enter any text message'});
            return;
        }

        var url = 'http://localhost:3001/api/v1/sms/send';
        var payload = {phone_number: data.get('number'),
            texts : [data.get('text1'),data.get('text2'),data.get('text2')]};

        fetch(url, {
            method: 'POST',
            body: JSON.stringify(payload),
            headers:{
                'Content-Type': 'application/json'
            }
        }).then(res => res.json())
            .then(data => this.setState({ error: data.message,result:data.status}))
            .catch(error => console.error('Error:', error));
     }

    render() {
        const {error,result} = this.state;
        const results = Object.entries(result).map(function([key, value])  {
            return  <li>Text {key} {value}</li>
        });
        return (
            <form onSubmit={this.handleSubmit}>
                <h1> SMS Sender </h1>
                <h2> Enter number and up to 3 texts to send! </h2>
                <p>{error}</p>
                <ul>{results}</ul>
                <div>
                    <label htmlFor="number">Phone number</label>
                    <input id="number" name="number" type="text" />
                </div>
                <div>
                   <label htmlFor="text1">Text 1</label>
                    <input id="text1" name="text1" type="text" />
                </div>
                <div>
                    <label htmlFor="text2">Text 2</label>
                    <input id="text2" name="text2" type="text" />
                </div>
                <div>
                    <label htmlFor="text3">Text 3</label>
                    <input id="text3" name="text3" type="text" />
                </div>
                <button>Send</button>
            </form>
        );
    }
}

export default hot(module)(App);
