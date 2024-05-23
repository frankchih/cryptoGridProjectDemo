import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"

const url = process.env.REACT_APP_URL

export const getLeverageAppApi = createApi({
    reducerPath: "getLeverageAppApi",
    baseQuery: fetchBaseQuery({ baseUrl: `${url}/api/leverageApp/` }),
    endpoints: (builder) => ({
        getLeverageSymbolList: builder.query({
            query: () => `getLeverageSymbolList/`,
        }),
        createLeverageSymbol: builder.mutation({
            query: (payload) => ({
                url: "createLeverageSymbol/",
                method: "POST",
                body: payload,
                headers: {
                    "Content-type": "application/json",
                },
            }),
            // invalidatesTags: ["Post"],
        }),
        updateLeverageSymbol: builder.mutation({
            query: (leverageSymbolId, payload) => ({
                url: `updateLeverageSymbol/${leverageSymbolId}/`,
                method: "PATCH",
                body: payload,
                headers: {
                    "Content-type": "application/json",
                },
            }),
        }),
        deleteLeverageSymbol: builder.mutation({
            query: (leverageSymbolId) => ({
                url: `deleteLeverageSymbol/${leverageSymbolId}/`,
                method: "DELETE",
                headers: {
                    "Content-type": "application/json",
                },
            }),
        }),
        getCurrAsset: builder.query({
            query: () => `getCurrAsset/`,
        }),
        postCurrAsset: builder.mutation({
            query: () => ({
                url: `getCurrAsset/`,
                method: "POST",
                headers: {
                    "Content-type": "application/json",
                },
            }),
        }),
    }),
})

export const {
    useGetLeverageSymbolListQuery,
    useCreateLeverageSymbolMutation,
    useUpdateLeverageSymbolMutation,
    useDeleteLeverageSymbolMutation,
    useGetCurrAssetQuery,
    usePostCurrAssetMutation
} = getLeverageAppApi
