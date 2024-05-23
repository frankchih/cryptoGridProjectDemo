import React from "react"
import styled from "styled-components"
import { TabView, TabPanel } from "primereact/tabview"

import QuoteDataTable from "../Quote/QuoteDataTable"
import LeverageSymbolDataTable from "../LeverageApp/LeverageSymbolDataTable"
import MainDescComponent from "../MainApp/MainDescComponent"
import ActivityLogDataTable from "../ActivityLog/ActivityLogDataTable"
import { Card } from "primereact/card"


const TempDiv = styled.div`
    height: 500px;
`

const MainContent = () => {
    return (
        <>
            <Card >
                此為 Demo Side Project，網格交易系統，策略大概是 預掛單買低賣高
                <br />
                有一些部分還在開發中，*已把下單相關程式刪除*，只是為了 demo 使用
            </Card>
            <br />
            <TabView>
                <TabPanel header="即時報價">
                    <div>背景有個 Goroutine 去呼叫幣安報價 websocket，並且送到 自己的websocket ，前端來接收即時值</div>
                    <QuoteDataTable />
                </TabPanel>
                <TabPanel header="網格設定">
                    <LeverageSymbolDataTable />
                </TabPanel>
                <TabPanel header="基本設定值">
                    <MainDescComponent />
                </TabPanel>
                <TabPanel header="Log紀錄">
                    <div>預期是要放歷程歷程記錄</div>
                    <ActivityLogDataTable />
                </TabPanel>
            </TabView>

            <TempDiv />
        </>
    )
}

export default MainContent
