import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"

const url = process.env.REACT_APP_URL

export const getOrderAppApi = createApi({
    reducerPath: "getOrderAppApi",
    baseQuery: fetchBaseQuery({ baseUrl: `${url}/api/orderApp/` }),
    endpoints: (builder) => ({
        createFirstGridOrder: builder.mutation({
            query: (payload) => ({
                url: "createFirstGridOrder/",
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
    useCreateFirstGridOrderMutation,
} = getOrderAppApi
