import React, { useEffect } from "react"

import { useGetSysCurrStatusQuery, usePostTaskRestartMutation } from "../../services/mainAppRTKService"
import { Button } from "primereact/button"

const MainDescShowComponent = ({ data }) => {
    return (
        <>
            <div>即時報價(taskQuoteValue): {data.taskQuoteValue}</div>
            <div>成交回報(taskOrderValue): {data.taskOrderValue}</div>
        </>
    )
}

const TaskRestartButton = ({ label, taskName, sysCurrStatusRefetch }) => {
    const [postTaskRestart, { isLoading, isError, error, isSuccess }] = usePostTaskRestartMutation()

    const postTaskRestartFunc = async () => {
        const formData = { taskName }
        await postTaskRestart(formData)
    }
    const handleBtnClick = (e) => {
        postTaskRestartFunc()
    }

    useEffect(() => {
        if (isSuccess) {
            //   toast.success('Post created successfully');
            console.log("新增成功")
            sysCurrStatusRefetch()
        }

        if (isError) {
            console.log("新增失敗")
        }
    }, [isLoading])

    return <Button label={label} onClick={() => handleBtnClick()} loading={isLoading} size="small"  style={{marginRight: "10px"}} />
}

const MainDescComponent = () => {
    const { data, error, isFetching, refetch } = useGetSysCurrStatusQuery()
    // console.log(data)

    return (
        <>
            <Button label="抓Redis當前資料" loading={isFetching} onClick={() => refetch()} size="small" style={{marginRight: "10px"}}/>
            <TaskRestartButton label="更新TaskQuote key" taskName={"TaskQuote"} sysCurrStatusRefetch={refetch} />
            <TaskRestartButton label="更新TaskOrder key" taskName={"TaskOrder"} sysCurrStatusRefetch={refetch} />
            <br /> 
            <br />
            兩個 Goroutine 狀態
            <br />
            {data && <MainDescShowComponent data={data} />}

        </>
    )
}

export default MainDescComponent
