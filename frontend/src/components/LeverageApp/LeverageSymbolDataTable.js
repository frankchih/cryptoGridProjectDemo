import React, { useEffect, useState } from "react"
import { Column } from "primereact/column"
import { DataTable } from "primereact/datatable"
import { Button } from "primereact/button"
import { Dialog } from "primereact/dialog"
import { InputText } from "primereact/inputtext"
import { ConfirmDialog, confirmDialog } from "primereact/confirmdialog"
import { Checkbox } from "primereact/checkbox"

import { useDispatch, useSelector } from "react-redux"

import { wsActions } from "../../slices/wsSlice"

import {
    getLeverageAppApi,
    useCreateLeverageSymbolMutation,
    useDeleteLeverageSymbolMutation,
    useGetCurrAssetQuery,
    useGetLeverageSymbolListQuery,
    usePostCurrAssetMutation,
} from "../../services/leverageAppRTKService"

import { useCreateFirstGridOrderMutation } from "../../services/orderAppRTKService"

const dateToText = (date) => {
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
}

const QuoteWsPrice = ({ symbol }) => {
    const quoteSymbolDict = useSelector((state) => state.quoteSlice.quoteSymbolDict)
    if (!(symbol in quoteSymbolDict)) {
        return ""
    }
    return quoteSymbolDict[symbol]["price"]
}
const QuoteWsEndTime = ({ symbol }) => {
    const quoteSymbolDict = useSelector((state) => state.quoteSlice.quoteSymbolDict)
    if (!(symbol in quoteSymbolDict)) {
        return ""
    }
    return dateToText(new Date(parseInt(quoteSymbolDict[symbol]["endTime"])))
}

const LeverageSymbolCreateButton = ({ leverageSymbolListRefetch }) => {
    const [visible, setVisible] = useState(false)
    const [symbol, setSymbol] = useState("")

    const [createLeverageSymbol, { isLoading, isError, error, isSuccess }] = useCreateLeverageSymbolMutation()

    const createLeverageSymbolFunc = async () => {
        if (symbol) {
            const formData = { symbol }
            await createLeverageSymbol(formData)
        }
    }
    const handleFormSubmit = (e) => {
        e.preventDefault()
        createLeverageSymbolFunc()
    }

    useEffect(() => {
        if (isSuccess) {
            //   toast.success('Post created successfully');
            console.log("新增成功")
            setVisible(false)
            leverageSymbolListRefetch()
        }

        if (isError) {
            console.log("新增失敗")
        }
    }, [isLoading])

    const footerContent = (
        <div>
            <Button label="關閉" icon="pi pi-check" onClick={() => setVisible(false)} size="small" />
        </div>
    )

    return (
        <>
            <Button label="新增" icon="pi pi-external-link" onClick={() => setVisible(true)} size="small" style={{ marginRight: "10px" }} />
            <Dialog header="Header" visible={visible} style={{ width: "50vw" }} onHide={() => setVisible(false)} footer={footerContent}>
                <form onSubmit={(e) => handleFormSubmit(e)}>
                    <div className="flex flex-column gap-2">
                        <label htmlFor="symbol">標的</label>
                        <InputText id="symbol" className="p-inputtext-sm" value={symbol} onChange={(e) => setSymbol(e.target.value)} />
                    </div>
                    <Button label="新增" icon="pi pi-check" loading={isLoading} size="small" />
                </form>
            </Dialog>
        </>
    )
}

const GetCurrAssetButton = ({ leverageSymbolListRefetch }) => {
    // const [trigger, { isLoading, isError, data, error }] = getLeverageAppApi.endpoints.getCurrAsset.useLazyQuery()
    // const { data, error, isLoading, refetch } = useGetCurrAssetQuery()

    // useEffect(() => {
    //     /* 這裡有問題 data 沒有變更 */
    //     console.log("GetCurrAssetButton useeffect", data)
    //     if (data) {
    //         // leverageSymbolListRefetch()
    //     }
    // }, [data])

    const [postCurrAsset, { isLoading, isError, error, isSuccess }] = usePostCurrAssetMutation()
    const postCurrAssetFunc = async () => {
        await postCurrAsset()
    }
    useEffect(() => {
        if (isSuccess) {
            leverageSymbolListRefetch()
            console.log("成功")
        }

        if (isError) {
            console.log("失敗")
        }
    }, [isLoading])

    const hanldeBtnClick = () => {
        // trigger()
        // refetch()
        postCurrAssetFunc()
    }

    return <Button label="更新當前倉位及交易對資訊" icon="pi pi-external-link" onClick={() => hanldeBtnClick()} isLoading={isLoading} size="small" style={{ marginRight: "10px" }} />
}

const LeverageSymbolDeleteButton = ({ leverageSymbolId, leverageSymbolListRefetch }) => {
    // const toast = useRef(null)
    const [deleteLeverageSymbol, { isLoading, isError, error, isSuccess }] = useDeleteLeverageSymbolMutation()
    const deleteLeverageSymbolFunc = async () => {
        await deleteLeverageSymbol(leverageSymbolId)
    }
    useEffect(() => {
        if (isSuccess) {
            leverageSymbolListRefetch()
            console.log("刪除成功")
        }

        if (isError) {
            console.log("刪除失敗")
        }
    }, [isLoading])

    const accept = () => {
        deleteLeverageSymbolFunc()
        // toast.current.show({
        //     severity: "info",
        //     summary: "Confirmed",
        //     detail: "You have accepted",
        //     life: 3000,
        // })
    }

    const reject = () => {
        // toast.current.show({
        //     severity: "warn",
        //     summary: "Rejected",
        //     detail: "You have rejected",
        //     life: 3000,
        // })
    }

    const confirm = () => {
        confirmDialog({
            message: `Do you want to delete this record? ${leverageSymbolId}`,
            header: "Delete Confirmation",
            icon: "pi pi-info-circle",
            acceptClassName: "p-button-danger",
            accept,
            reject,
        })
    }

    return (
        <>
            {/* <Toast ref={toast} /> */}

            <div className="card flex flex-wrap gap-2 justify-content-center">
                <Button onClick={confirm} icon="pi pi-times" label="刪除" style={{ padding: "0 10px" }} />
            </div>
        </>
    )
}

const OrderSymbolDownDataTable = ({ orderSymbolDownList }) => {
    // const { data, error, isLoading, refetch } = useGetLeverageSymbolListQuery()
    // const quoteDataList = useSelector((state) => state.quoteSlice.quoteDataList)
    // const dispatch = useDispatch()
    useEffect(() => {
        // dispatch(wsActions.startConnecting())
        return () => {
            // dispatch({ type: 'socket/disconnect' })
        }
    }, [])

    const refreshBtnClick = () => {
        // refetch()
    }

    return (
        <DataTable
            header={"網格下單"}
            value={orderSymbolDownList}
            size={"small"}
            // style={{ textAlign: "center" }}
            scrollable
            scrollHeight="400px"
            showGridlines
        >
            <Column field="ID" header="ID" style={{ minWidth: "100px", padding: 0 }}></Column>
            <Column field="symbol" header="交易對" style={{ minWidth: "100px", padding: 0 }}></Column>
            <Column field="price" header="價" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
            <Column field="executedQty" header="數量" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
        </DataTable>
    )
}

const OrderSymbolUpDataTable = ({ orderSymbolUpList }) => {
    // const { data, error, isLoading, refetch } = useGetLeverageSymbolListQuery()

    // const quoteDataList = useSelector((state) => state.quoteSlice.quoteDataList)

    // const dispatch = useDispatch()
    useEffect(() => {
        // dispatch(wsActions.startConnecting())
        return () => {
            // dispatch({ type: 'socket/disconnect' })
        }
    }, [])

    const refreshBtnClick = () => {
        // refetch()
    }

    return (
        <DataTable
            header={"網格上單"}
            value={orderSymbolUpList}
            size={"small"}
            // style={{ textAlign: "center" }}
            scrollable
            scrollHeight="400px"
            showGridlines
        >
            <Column field="ID" header="ID" style={{ minWidth: "100px", padding: 0 }}></Column>
            <Column field="symbol" header="交易對" style={{ minWidth: "100px", padding: 0 }}></Column>
            <Column field="price" header="價" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
            <Column field="executedQty" header="數量" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
        </DataTable>
    )
}

const OrderSymbolCreateFirstGridOrderButton = ({ leverageSymbolId, leverageSymbolListRefetch, rowData }) => {
    const [visible, setVisible] = useState(false)
    const [symbol, setSymbol] = useState(rowData.symbol)

    const [settingPriceStepPercent, setSettingPriceStepPercent] = useState("0.03")
    const [settingCalcNum, setSettingCalcNum] = useState(10)
    const [inventoryVolume, setInventoryVolume] = useState(rowData.inventoryVolume)
    const [isSimulate, setIsSimulate] = useState(true)

    const quoteSymbolDict = useSelector((state) => state.quoteSlice.quoteSymbolDict)

    const [marketPrice, setMarketPrice] = useState("0")
    const [orderSymbolDownList, setOrderSymbolDownList] = useState(null)
    const [orderSymbolUpList, setOrderSymbolUpList] = useState(null)

    useEffect(() => {
        if (rowData.symbol in quoteSymbolDict) {
            const _marketPrice = quoteSymbolDict[rowData.symbol]["price"]
            setMarketPrice(_marketPrice)
        }
    }, [quoteSymbolDict])

    const [createFirstGridOrder, result] = useCreateFirstGridOrderMutation()

    const createFirstGridOrderFunc = async () => {
        const formData = {
            leverageSymbolId,
            settingPriceStepPercent,
            settingCalcNum: parseInt(settingCalcNum),
            inventoryVolume,
            marketPrice,
            isSimulate,
        }

        const res = await createFirstGridOrder(formData).unwrap()
        if (res.msg) {
            setOrderSymbolDownList(res.orderSymbolDownList)
            setOrderSymbolUpList(res.orderSymbolUpList)
        }
    }
    const handleFormSubmit = (e) => {
        e.preventDefault()
        createFirstGridOrderFunc()
    }

    // useEffect(() => {
    //     console.log(result)
    // }, [result])

    // useEffect(() => {
    //     if (isSuccess) {
    //         //   toast.success('Post created successfully');
    //         console.log("新增成功")
    //         setVisible(false)
    //         leverageSymbolListRefetch()
    //     }
    // }, [isLoading])

    // useEffect(() => {
    //     if (isError) {
    //         console.log(isError, error)
    //     }
    // }, [isError])

    // useEffect(() => {
    //     if (isSuccess) {
    //         //   toast.success('Post created successfully');
    //         console.log("新增成功")
    //         // setVisible(false)
    //         leverageSymbolListRefetch()
    //     }
    // }, [isSuccess])
    // console.log(isLoading, isError, error, isSuccess, data)

    const footerContent = (
        <div>
            <Button label="關閉" icon="pi pi-check" onClick={() => setVisible(false)} size="small" />
        </div>
    )

    return (
        <>
            <Button label="模擬單計算" icon="pi pi-external-link" onClick={() => setVisible(true)} size="small" />
            <Dialog header="Header" visible={visible} style={{ width: "50vw" }} onHide={() => setVisible(false)} footer={footerContent} blockScroll>
                <form onSubmit={(e) => handleFormSubmit(e)}>
                    <div className="flex flex-column gap-2">
                        <label htmlFor="symbol">標的</label>
                        <InputText id="symbol" className="p-inputtext-sm" value={symbol} onChange={(e) => setSymbol(e.target.value)} />
                    </div>

                    <div className="flex flex-column gap-2">
                        <label htmlFor="settingPriceStepPercent">設定上下價格%</label>
                        <InputText
                            id="settingPriceStepPercent"
                            className="p-inputtext-sm"
                            value={settingPriceStepPercent}
                            onChange={(e) => setSettingPriceStepPercent(e.target.value)}
                        />
                    </div>
                    <div className="flex flex-column gap-2">
                        <label htmlFor="settingCalcNum">設定 上下筆數</label>
                        <InputText
                            id="settingCalcNum"
                            className="p-inputtext-sm"
                            value={settingCalcNum}
                            onChange={(e) => setSettingCalcNum(e.target.value)}
                        />
                    </div>
                    
                    <div className="flex flex-column gap-2">
                        <label htmlFor="inventoryVolume">庫存量</label>
                        <InputText
                            id="inventoryVolume"
                            className="p-inputtext-sm"
                            value={inventoryVolume}
                            onChange={(e) => setInventoryVolume(e.target.value)}
                        />
                    </div>
                    <div className="flex flex-column gap-2">
                        <label htmlFor="marketPrice">市價</label>
                        <InputText id="marketPrice" className="p-inputtext-sm" value={marketPrice} onChange={(e) => setMarketPrice(e.target.value)} />
                    </div>
                    <div className="flex flex-column gap-2">
                        <label htmlFor="isSimulate">模擬單</label>
                        {/* <Checkbox onChange={(e) => setIsSimulate(e.checked)} checked={isSimulate}></Checkbox> */}
                    </div>

                    <Button
                        label="模擬預期的下單(不會真的下單)"
                        icon="pi pi-check"
                        // loading={isLoading}
                        size="small"
                    />
                </form>

                <OrderSymbolDownDataTable orderSymbolDownList={orderSymbolDownList} />
                <OrderSymbolUpDataTable orderSymbolUpList={orderSymbolUpList} />
            </Dialog>
        </>
    )
}

const LeverageSymbolDataTable = () => {
    const { data, error, isLoading, refetch } = useGetLeverageSymbolListQuery()

    const refreshBtnClick = () => {
        refetch()
    }

    return (
        <div>
            <ConfirmDialog />
            <LeverageSymbolCreateButton leverageSymbolListRefetch={refetch} />
            <GetCurrAssetButton leverageSymbolListRefetch={refetch} />
            <Button label="重新整理" onClick={refreshBtnClick} loading={isLoading} size="small" />
            <DataTable
                header={"槓桿"}
                loading={isLoading}
                value={data?.leverageSymbols}
                size={"small"}
                // style={{ textAlign: "center" }}
                scrollable
                scrollHeight="500px"
                showGridlines
            >
                <Column
                    field="options"
                    header="刪除"
                    body={(rowData) => {
                        return (
                            <>
                                <LeverageSymbolDeleteButton leverageSymbolId={rowData.ID} leverageSymbolListRefetch={refetch} />
                            </>
                        )
                    }}
                    style={{ minWidth: "80px", textAlign: "center" }}    
                ></Column>
                <Column
                    field="options"
                    header="功能"
                    body={(rowData) => {
                        return (
                            <>
                                <OrderSymbolCreateFirstGridOrderButton
                                    leverageSymbolId={rowData.ID}
                                    leverageSymbolListRefetch={refetch}
                                    rowData={rowData}
                                />
                            </>
                        )
                    }}
                    style={{ minWidth: "150px", textAlign: "center" }}    
                ></Column>

                <Column field="symbol" header="交易對" style={{ minWidth: "100px", padding: 5 }}></Column>
                <Column
                    field="marketPrice"
                    header="市價"
                    style={{ minWidth: "50px", padding: 5, textAlign: "center" }}
                    body={(rowData) => {
                        return <QuoteWsPrice symbol={rowData.symbol} />
                    }}
                ></Column>
                <Column
                    field="endTime"
                    header="市價更新時間"
                    body={(rowData) => {
                        return <QuoteWsEndTime symbol={rowData.symbol} />
                    }}
                ></Column>
                <Column field="inventoryVolume" header="庫存" style={{ minWidth: "50px", padding: 5, textAlign: "center" }}></Column>
                <Column
                    field="inventoryVolumeUpdateDateTime"
                    header="庫存更新時間"
                    style={{ minWidth: "50px", padding: 5, textAlign: "center" }}
                    body={(rowData) => {
                        if (rowData.inventoryVolumeUpdateDateTime !== "0001-01-01T08:00:00+08:00") {
                            const date = new Date(rowData.inventoryVolumeUpdateDateTime)
                            const hours = date.getHours().toString().padStart(2, "0")
                            const minutes = date.getMinutes().toString().padStart(2, "0")
                            const seconds = date.getSeconds().toString().padStart(2, "0")
                            const timeString = `${hours}:${minutes}:${seconds}`

                            return timeString
                        }
                        return ""
                    }}
                ></Column>

                <Column field="pricePrecision" header="pricePrecision" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
                <Column field="quantityPrecision" header="quantityPrecision" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
                <Column field="baseAssetPrecision" header="baseAssetPrecision" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
                <Column field="quotePrecision" header="quotePrecision" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
                <Column field="filterMinNotional" header="filterMinNotional" style={{ minWidth: "50px", padding: 0, textAlign: "center" }}></Column>
                <Column
                    field="isSimulate"
                    header="isSimulate"
                    style={{ minWidth: "50px", padding: 0, textAlign: "center" }}
                    body={(rowData) => rowData.isSimulate}
                ></Column>

                
            </DataTable>
        </div>
    )
}

export default LeverageSymbolDataTable
