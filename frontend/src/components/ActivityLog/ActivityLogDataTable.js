import React, { useState } from "react"
import styled from "styled-components"
import { DataTable } from "primereact/datatable"
import { Column } from "primereact/column"
import { Button } from "primereact/button"
import { useDispatch, useSelector } from "react-redux"
import { useEffect } from "react"
import { activityLogGetListAsync } from "../../slices/activityLogSlice"



const ActivityLogDataTable = () => {
    const activityLogs = useSelector((state) => state.activityLogReducer.activityLogs)
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(activityLogGetListAsync())
    }, [])

    return (
        <>
            <Button label="紀錄" size="small" onClick={() => dispatch(activityLogGetListAsync())} loading={!activityLogs} />
            <DataTable value={activityLogs} size={"small"} style={{ minWidth: "50rem" }} scrollable scrollHeight="400px">
                <Column field="message" header="訊息"></Column>
                <Column field="CreatedAt" header="建立時間"></Column>
            </DataTable>
        </>
    )
}

export default ActivityLogDataTable
