import React from "react"
import styled from "styled-components"
import MainContent from "./MainContent"

const GridContainer = styled.div`
    width: 100%;
    display: flex;
`

const LeftItem = styled.div`
    width: 200px; 
    background-color: #ccc;
    padding: 20px;
`

const RightItem = styled.div`
    width: calc(100% - 200px);
    background-color: #eee;
    padding: 20px;
`

const Layout = () => {
    return (
        <GridContainer>
            <LeftItem>left</LeftItem>
            <RightItem>
                <MainContent />
            </RightItem>
        </GridContainer>
    )
}

export default Layout
