const url = process.env.REACT_APP_URL

class ActivityLogServie {
    
    static async getList() {
        const response = await fetch(`${url}/api/activityLog/getList/`)
        const responseJson = await response.json()
        return {response, responseJson}
    }
}
export default ActivityLogServie