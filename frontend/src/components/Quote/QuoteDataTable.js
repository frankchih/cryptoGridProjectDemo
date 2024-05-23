import React, { useEffect } from "react"
import { Column } from "primereact/column"
import { DataTable } from "primereact/datatable"
import { useDispatch, useSelector } from "react-redux"
import { wsActions } from "../../slices/wsSlice"
const dateToText = (date) => {
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
}
const QuoteDataTable = () => {
    const quoteDataList = useSelector(
        (state) => state.quoteSlice.quoteDataList
    )
    const dispatch = useDispatch()
    useEffect(() => {

        dispatch(wsActions.startConnecting())
        return () => {
            // dispatch({ type: 'socket/disconnect' })
        }
    }, [])
    return (
        <DataTable
            value={quoteDataList}
            size={"small"}
            style={{ minWidth: "50rem" }}
            scrollable
            scrollHeight="400px"
        >
            <Column field="symbol" header="交易對"></Column>
            <Column field="price" header="目前價格"></Column>
            <Column field="endTime" header="最後更新時間" body={(rowData) => {
                return dateToText(new Date(parseInt(rowData["endTime"])))
            }}></Column>
        </DataTable>
    )
}

export default QuoteDataTable
