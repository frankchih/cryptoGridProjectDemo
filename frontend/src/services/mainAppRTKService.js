import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"

const url = process.env.REACT_APP_URL

export const getMainAppApi = createApi({
    reducerPath: "getMainAppApi",
    baseQuery: fetchBaseQuery({ baseUrl: `${url}/api/mainApp/` }),
    endpoints: (builder) => ({
        getSysCurrStatus: builder.query({
            query: () => `getSysCurrStatus/`,
        }),
        postTaskRestart: builder.mutation({
            query: (payload) => ({
                url: "postTaskRestart/",
                method: "POST",
                body: payload,
                headers: {
                    "Content-type": "application/json",
                },
            }),
            // invalidatesTags: ["Post"],
        }),
    }),
})

export const {
    useGetSysCurrStatusQuery,
    usePostTaskRestartMutation,
} = getMainAppApi
