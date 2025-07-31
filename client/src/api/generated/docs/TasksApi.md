# TasksApi

All URIs are relative to *http://localhost:4000/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**tasksGet**](#tasksget) | **GET** /tasks | Get all tasks with filtering options|
|[**tasksPost**](#taskspost) | **POST** /tasks | Create a new task|
|[**tasksTaskIdGet**](#taskstaskidget) | **GET** /tasks/{task_id} | Get a task by ID|

# **tasksGet**
> TasksGet200Response tasksGet()

Retrieves tasks based on filter criteria (active/archived/completed) and detail level

### Example

```typescript
import {
    TasksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let authorization: string; //Bearer token (default to undefined)
let detailLevel: 'minimal' | 'standard' | 'full'; //Detail level (optional) (default to undefined)
let showActive: boolean; //Include active tasks (default: true) (optional) (default to undefined)
let showArchived: boolean; //Include archived tasks (default: false) (optional) (default to undefined)
let showCompleted: boolean; //Include completed tasks (default: true) (optional) (default to undefined)

const { status, data } = await apiInstance.tasksGet(
    authorization,
    detailLevel,
    showActive,
    showArchived,
    showCompleted
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **authorization** | [**string**] | Bearer token | defaults to undefined|
| **detailLevel** | [**&#39;minimal&#39; | &#39;standard&#39; | &#39;full&#39;**]**Array<&#39;minimal&#39; &#124; &#39;standard&#39; &#124; &#39;full&#39;>** | Detail level | (optional) defaults to undefined|
| **showActive** | [**boolean**] | Include active tasks (default: true) | (optional) defaults to undefined|
| **showArchived** | [**boolean**] | Include archived tasks (default: false) | (optional) defaults to undefined|
| **showCompleted** | [**boolean**] | Include completed tasks (default: true) | (optional) defaults to undefined|


### Return type

**TasksGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of tasks |  -  |
|**400** | Invalid query parameters |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **tasksPost**
> TasksPost201Response tasksPost(createTaskRequest)

Create a new task for the authenticated user

### Example

```typescript
import {
    TasksApi,
    Configuration,
    DtoCreateTaskRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let authorization: string; //Bearer token (default to undefined)
let createTaskRequest: DtoCreateTaskRequest; //Task creation data

const { status, data } = await apiInstance.tasksPost(
    authorization,
    createTaskRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createTaskRequest** | **DtoCreateTaskRequest**| Task creation data | |
| **authorization** | [**string**] | Bearer token | defaults to undefined|


### Return type

**TasksPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Task successfully created |  -  |
|**400** | Invalid request format |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **tasksTaskIdGet**
> TasksTaskIdGet200Response tasksTaskIdGet()

Retrieves a single task by its unique identifier

### Example

```typescript
import {
    TasksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let authorization: string; //Bearer token (default to undefined)
let taskId: string; //Task ID (default to undefined)

const { status, data } = await apiInstance.tasksTaskIdGet(
    authorization,
    taskId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **authorization** | [**string**] | Bearer token | defaults to undefined|
| **taskId** | [**string**] | Task ID | defaults to undefined|


### Return type

**TasksTaskIdGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Task details |  -  |
|**400** | Invalid task ID format |  -  |
|**401** | Unauthorized |  -  |
|**404** | Task not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

