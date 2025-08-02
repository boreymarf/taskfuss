# ApiError

Common error response structure for API failures

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **string** |  | [optional] [default to undefined]
**details** | **object** |  | [optional] [default to undefined]
**message** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { ApiError } from './api';

const instance: ApiError = {
    code,
    details,
    message,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
