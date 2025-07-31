# AuthenticationApi

All URIs are relative to *http://localhost:4000/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authLoginPost**](#authloginpost) | **POST** /auth/login | User login|
|[**authRegisterPost**](#authregisterpost) | **POST** /auth/register | Register a new user|

# **authLoginPost**
> AuthLoginPost200Response authLoginPost(loginRequest)

Authenticate user credentials and return a JWT token

### Example

```typescript
import {
    AuthenticationApi,
    Configuration,
    DtoLoginRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthenticationApi(configuration);

let loginRequest: DtoLoginRequest; //Login credentials

const { status, data } = await apiInstance.authLoginPost(
    loginRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **loginRequest** | **DtoLoginRequest**| Login credentials | |


### Return type

**AuthLoginPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully authenticated |  -  |
|**400** | Invalid request format |  -  |
|**401** | Invalid credentials |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRegisterPost**
> AuthRegisterPost201Response authRegisterPost(registerRequest)

Create a new user account and return a JWT token

### Example

```typescript
import {
    AuthenticationApi,
    Configuration,
    DtoRegisterRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthenticationApi(configuration);

let registerRequest: DtoRegisterRequest; //User registration data

const { status, data } = await apiInstance.authRegisterPost(
    registerRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **registerRequest** | **DtoRegisterRequest**| User registration data | |


### Return type

**AuthRegisterPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Successfully registered |  -  |
|**400** | Invalid request format (code: BAD_REQUEST) or username/email already exists (code: DUPLICATE_USER) |  -  |
|**500** | Internal server error (code: INTERNAL_ERROR) |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

