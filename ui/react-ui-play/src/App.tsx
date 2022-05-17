import React from 'react';
import logo from './logo.svg';
import './App.css';

import { Refine } from "@pankod/refine-core";
import { Layout, ReadyPage, notificationProvider, ErrorComponent } from "@pankod/refine-antd";
import routerProvider from "@pankod/refine-react-router-v6";
import dataProvider from "@pankod/refine-simple-rest";

import "@pankod/refine-antd/dist/styles.min.css";

import { RaceList, RaceShow } from "./pages/races";

function App() {
  return (
    <div className="App">
      <Refine
        routerProvider={routerProvider}
        dataProvider={dataProvider("http://localhost:8080")}
        Layout={Layout}
        ReadyPage={ReadyPage}
        notificationProvider={notificationProvider}
        catchAll={<ErrorComponent />}
        resources={[{ name: "races", list: RaceList, show: RaceShow}]}
      />
    </div>
  );
}

export default App;
