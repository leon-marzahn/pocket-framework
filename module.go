package pocketframework

import (
  "github.com/pocketbase/pocketbase/core"
  "github.com/pocketbase/pocketbase/tools/hook"
)

type Module interface {
  // Prefix should return the api path prefix for this module. This might be nested depending on if this module has a parent.
  Prefix() string

  // RegisterHooks can be used to register hooks in pocketbase.
  RegisterHooks(app ModuleAppHooks) error

  // RegisterRoutes can be used to register routes in pocketbase.
  RegisterRoutes(groups RouterGroups) error
}

type ModuleWithChildren interface {
  Module
  Children() []Module
}

type ModuleAppHooks interface {
  // OnBootstrap hook is triggered when initializing the main application
  // resources (db, app settings, etc).
  OnBootstrap() *hook.Hook[*core.BootstrapEvent]

  // OnTerminate hook is triggered when the app is in the process
  // of being terminated (ex. on SIGTERM signal).
  //
  // Note that the app could be terminated abruptly without awaiting the hook completion.
  OnTerminate() *hook.Hook[*core.TerminateEvent]

  // OnBackupCreate hook is triggered on each [App.CreateBackup] call.
  OnBackupCreate() *hook.Hook[*core.BackupEvent]

  // OnBackupRestore hook is triggered before app backup restore (aka. [App.RestoreBackup] call).
  //
  // Note that by default on success the application is restarted and the after state of the hook is ignored.
  OnBackupRestore() *hook.Hook[*core.BackupEvent]

  // ---------------------------------------------------------------
  // DB models event hooks
  // ---------------------------------------------------------------

  // OnModelValidate is triggered every time when a model is being validated
  // (e.g. triggered by App.Validate() or App.Save()).
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelValidate(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // ---------------------------------------------------------------

  // OnModelCreate is triggered every time when a new model is being created
  // (e.g. triggered by App.Save()).
  //
  // Operations BEFORE the e.Next() execute before the model validation
  // and the INSERT DB statement.
  //
  // Operations AFTER the e.Next() execute after the model validation
  // and the INSERT DB statement.
  //
  // Note that successful execution doesn't guarantee that the model
  // is persisted in the database since its wrapping transaction may
  // not have been committed yet.
  // If you want to listen to only the actual persisted events, you can
  // bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelCreate(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelCreateExecute is triggered after successful Model validation
  // and right before the model INSERT DB statement execution.
  //
  // Usually it is triggered as part of the App.Save() in the following firing order:
  // OnModelCreate {
  //    -> OnModelValidate (skipped with App.SaveNoValidate())
  //    -> OnModelCreateExecute
  // }
  //
  // Note that successful execution doesn't guarantee that the model
  // is persisted in the database since its wrapping transaction may have been
  // committed yet.
  // If you want to listen to only the actual persisted events,
  // you can bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelCreateExecute(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterCreateSuccess is triggered after each successful
  // Model DB create persistence.
  //
  // Note that when a Model is persisted as part of a transaction,
  // this hook is delayed and executed only AFTER the transaction has been committed.
  // This hook is NOT triggered in case the transaction rollbacks
  // (aka. when the model wasn't persisted).
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterCreateSuccess(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterCreateError is triggered after each failed
  // Model DB create persistence.
  //
  // Note that the execution of this hook is either immediate or delayed
  // depending on the error:
  //   - "immediate" on App.Save() failure
  //   - "delayed" on transaction rollback
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterCreateError(tags ...string) *hook.TaggedHook[*core.ModelErrorEvent]

  // ---------------------------------------------------------------

  // OnModelUpdate is triggered every time when a new model is being updated
  // (e.g. triggered by App.Save()).
  //
  // Operations BEFORE the e.Next() execute before the model validation
  // and the UPDATE DB statement.
  //
  // Operations AFTER the e.Next() execute after the model validation
  // and the UPDATE DB statement.
  //
  // Note that successful execution doesn't guarantee that the model
  // is persisted in the database since its wrapping transaction may
  // not have been committed yet.
  // If you want to listen to only the actual persisted events, you can
  // bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelUpdate(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelUpdateExecute is triggered after successful Model validation
  // and right before the model UPDATE DB statement execution.
  //
  // Usually it is triggered as part of the App.Save() in the following firing order:
  // OnModelUpdate {
  //    -> OnModelValidate (skipped with App.SaveNoValidate())
  //    -> OnModelUpdateExecute
  // }
  //
  // Note that successful execution doesn't guarantee that the model
  // is persisted in the database since its wrapping transaction may have been
  // committed yet.
  // If you want to listen to only the actual persisted events,
  // you can bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelUpdateExecute(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterUpdateSuccess is triggered after each successful
  // Model DB update persistence.
  //
  // Note that when a Model is persisted as part of a transaction,
  // this hook is delayed and executed only AFTER the transaction has been committed.
  // This hook is NOT triggered in case the transaction rollbacks
  // (aka. when the model changes weren't persisted).
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterUpdateError is triggered after each failed
  // Model DB update persistence.
  //
  // Note that the execution of this hook is either immediate or delayed
  // depending on the error:
  //   - "immediate" on App.Save() failure
  //   - "delayed" on transaction rollback
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterUpdateError(tags ...string) *hook.TaggedHook[*core.ModelErrorEvent]

  // ---------------------------------------------------------------

  // OnModelDelete is triggered every time when a new model is being deleted
  // (e.g. triggered by App.Delete()).
  //
  // Note that successful execution doesn't guarantee that the model
  // is deleted from the database since its wrapping transaction may
  // not have been committed yet.
  // If you want to listen to only the actual persisted deleted events, you can
  // bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelDelete(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelUpdateExecute is triggered right before the model
  // DELETE DB statement execution.
  //
  // Usually it is triggered as part of the App.Delete() in the following firing order:
  // OnModelDelete {
  //    -> (internal delete checks)
  //    -> OnModelDeleteExecute
  // }
  //
  // Note that successful execution doesn't guarantee that the model
  // is deleted from the database since its wrapping transaction may
  // not have been committed yet.
  // If you want to listen to only the actual persisted deleted events, you can
  // bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelDeleteExecute(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterDeleteSuccess is triggered after each successful
  // Model DB delete persistence.
  //
  // Note that when a Model is deleted as part of a transaction,
  // this hook is delayed and executed only AFTER the transaction has been committed.
  // This hook is NOT triggered in case the transaction rollbacks
  // (aka. when the model delete wasn't persisted).
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*core.ModelEvent]

  // OnModelAfterDeleteError is triggered after each failed
  // Model DB delete persistence.
  //
  // Note that the execution of this hook is either immediate or delayed
  // depending on the error:
  //   - "immediate" on App.Delete() failure
  //   - "delayed" on transaction rollback
  //
  // For convenience, if you want to listen to only the Record models
  // events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
  //
  // If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnModelAfterDeleteError(tags ...string) *hook.TaggedHook[*core.ModelErrorEvent]

  // ---------------------------------------------------------------
  // Record models event hooks
  // ---------------------------------------------------------------

  // OnRecordEnrich is triggered every time when a record is enriched
  // (as part of the builtin Record responses, during realtime message seriazation, or when [apis.EnrichRecord] is invoked).
  //
  // It could be used for example to redact/hide or add computed temporary
  // Record model props only for the specific request info. For example:
  //
  //  app.OnRecordEnrich("posts").BindFunc(func(e core.*RecordEnrichEvent) {
  //      // hide one or more fields
  //      e.Record.Hide("role")
  //
  //      // add new custom field for registered users
  //      if e.RequestInfo.Auth != nil && e.RequestInfo.Auth.Collection().Name == "users" {
  //          e.Record.WithCustomData(true) // for security requires explicitly allowing it
  //          e.Record.Set("computedScore", e.Record.GetInt("score") * e.RequestInfo.Auth.GetInt("baseScore"))
  //      }
  //
  //      return e.Next()
  //  })
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordEnrich(tags ...string) *hook.TaggedHook[*core.RecordEnrichEvent]

  // OnRecordValidate is a Record proxy model hook of [OnModelValidate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordValidate(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // ---------------------------------------------------------------

  // OnRecordCreate is a Record proxy model hook of [OnModelCreate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordCreate(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordCreateExecute is a Record proxy model hook of [OnModelCreateExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordCreateExecute(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterCreateSuccess is a Record proxy model hook of [OnModelAfterCreateSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterCreateSuccess(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterCreateError is a Record proxy model hook of [OnModelAfterCreateError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterCreateError(tags ...string) *hook.TaggedHook[*core.RecordErrorEvent]

  // ---------------------------------------------------------------

  // OnRecordUpdate is a Record proxy model hook of [OnModelUpdate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordUpdate(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordUpdateExecute is a Record proxy model hook of [OnModelUpdateExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordUpdateExecute(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterUpdateSuccess is a Record proxy model hook of [OnModelAfterUpdateSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterUpdateError is a Record proxy model hook of [OnModelAfterUpdateError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterUpdateError(tags ...string) *hook.TaggedHook[*core.RecordErrorEvent]

  // ---------------------------------------------------------------

  // OnRecordDelete is a Record proxy model hook of [OnModelDelete].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordDelete(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordDeleteExecute is a Record proxy model hook of [OnModelDeleteExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordDeleteExecute(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterDeleteSuccess is a Record proxy model hook of [OnModelAfterDeleteSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*core.RecordEvent]

  // OnRecordAfterDeleteError is a Record proxy model hook of [OnModelAfterDeleteError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAfterDeleteError(tags ...string) *hook.TaggedHook[*core.RecordErrorEvent]

  // ---------------------------------------------------------------
  // Collection models event hooks
  // ---------------------------------------------------------------

  // OnCollectionValidate is a Collection proxy model hook of [OnModelValidate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionValidate(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // ---------------------------------------------------------------

  // OnCollectionCreate is a Collection proxy model hook of [OnModelCreate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionCreate(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionCreateExecute is a Collection proxy model hook of [OnModelCreateExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionCreateExecute(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterCreateSuccess is a Collection proxy model hook of [OnModelAfterCreateSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterCreateSuccess(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterCreateError is a Collection proxy model hook of [OnModelAfterCreateError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterCreateError(tags ...string) *hook.TaggedHook[*core.CollectionErrorEvent]

  // ---------------------------------------------------------------

  // OnCollectionUpdate is a Collection proxy model hook of [OnModelUpdate].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionUpdate(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionUpdateExecute is a Collection proxy model hook of [OnModelUpdateExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionUpdateExecute(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterUpdateSuccess is a Collection proxy model hook of [OnModelAfterUpdateSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterUpdateError is a Collection proxy model hook of [OnModelAfterUpdateError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterUpdateError(tags ...string) *hook.TaggedHook[*core.CollectionErrorEvent]

  // ---------------------------------------------------------------

  // OnCollectionDelete is a Collection proxy model hook of [OnModelDelete].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionDelete(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionDeleteExecute is a Collection proxy model hook of [OnModelDeleteExecute].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionDeleteExecute(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterDeleteSuccess is a Collection proxy model hook of [OnModelAfterDeleteSuccess].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*core.CollectionEvent]

  // OnCollectionAfterDeleteError is a Collection proxy model hook of [OnModelAfterDeleteError].
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnCollectionAfterDeleteError(tags ...string) *hook.TaggedHook[*core.CollectionErrorEvent]

  // ---------------------------------------------------------------
  // Mailer event hooks
  // ---------------------------------------------------------------

  // OnMailerSend hook is triggered every time when a new email is
  // being send using the [App.NewMailClient()] instance.
  //
  // It allows intercepting the email message or to use a custom mailer client.
  OnMailerSend() *hook.Hook[*core.MailerEvent]

  // OnMailerRecordAuthAlertSend hook is triggered when
  // sending a new device login auth alert email, allowing you to
  // intercept and customize the email message that is being sent.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnMailerRecordAuthAlertSend(tags ...string) *hook.TaggedHook[*core.MailerRecordEvent]

  // OnMailerRecordPasswordResetSend hook is triggered when
  // sending a password reset email to an auth record, allowing
  // you to intercept and customize the email message that is being sent.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnMailerRecordPasswordResetSend(tags ...string) *hook.TaggedHook[*core.MailerRecordEvent]

  // OnMailerRecordVerificationSend hook is triggered when
  // sending a verification email to an auth record, allowing
  // you to intercept and customize the email message that is being sent.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnMailerRecordVerificationSend(tags ...string) *hook.TaggedHook[*core.MailerRecordEvent]

  // OnMailerRecordEmailChangeSend hook is triggered when sending a
  // confirmation new address email to an auth record, allowing
  // you to intercept and customize the email message that is being sent.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnMailerRecordEmailChangeSend(tags ...string) *hook.TaggedHook[*core.MailerRecordEvent]

  // OnMailerRecordOTPSend hook is triggered when sending an OTP email
  // to an auth record, allowing you to intercept and customize the
  // email message that is being sent.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnMailerRecordOTPSend(tags ...string) *hook.TaggedHook[*core.MailerRecordEvent]

  // ---------------------------------------------------------------
  // Realtime API event hooks
  // ---------------------------------------------------------------

  // OnRealtimeConnectRequest hook is triggered when establishing the SSE client connection.
  //
  // Any execution after e.Next() of a hook handler happens after the client disconnects.
  OnRealtimeConnectRequest() *hook.Hook[*core.RealtimeConnectRequestEvent]

  // OnRealtimeMessageSend hook is triggered when sending an SSE message to a client.
  OnRealtimeMessageSend() *hook.Hook[*core.RealtimeMessageEvent]

  // OnRealtimeSubscribeRequest hook is triggered when updating the
  // client subscriptions, allowing you to further validate and
  // modify the submitted change.
  OnRealtimeSubscribeRequest() *hook.Hook[*core.RealtimeSubscribeRequestEvent]

  // ---------------------------------------------------------------
  // Settings API event hooks
  // ---------------------------------------------------------------

  // OnSettingsListRequest hook is triggered on each API Settings list request.
  //
  // Could be used to validate or modify the response before returning it to the client.
  OnSettingsListRequest() *hook.Hook[*core.SettingsListRequestEvent]

  // OnSettingsUpdateRequest hook is triggered on each API Settings update request.
  //
  // Could be used to additionally validate the request data or
  // implement completely different persistence behavior.
  OnSettingsUpdateRequest() *hook.Hook[*core.SettingsUpdateRequestEvent]

  // OnSettingsReload hook is triggered every time when the App.Settings()
  // is being replaced with a new state.
  //
  // Calling App.Settings() after e.Next() returns the new state.
  OnSettingsReload() *hook.Hook[*core.SettingsReloadEvent]

  // ---------------------------------------------------------------
  // File API event hooks
  // ---------------------------------------------------------------

  // OnFileDownloadRequest hook is triggered before each API File download request.
  //
  // Could be used to validate or modify the file response before
  // returning it to the client.
  OnFileDownloadRequest(tags ...string) *hook.TaggedHook[*core.FileDownloadRequestEvent]

  // OnFileBeforeTokenRequest hook is triggered on each auth file token API request.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnFileTokenRequest(tags ...string) *hook.TaggedHook[*core.FileTokenRequestEvent]

  // ---------------------------------------------------------------
  // Record Auth API event hooks
  // ---------------------------------------------------------------

  // OnRecordAuthRequest hook is triggered on each successful API
  // record authentication request (sign-in, token refresh, etc.).
  //
  // Could be used to additionally validate or modify the authenticated
  // record data and token.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAuthRequest(tags ...string) *hook.TaggedHook[*core.RecordAuthRequestEvent]

  // OnRecordAuthWithPasswordRequest hook is triggered on each
  // Record auth with password API request.
  //
  // [RecordAuthWithPasswordRequestEvent.Record] could be nil if no matching identity is found, allowing
  // you to manually locate a different Record model (by reassigning [RecordAuthWithPasswordRequestEvent.Record]).
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*core.RecordAuthWithPasswordRequestEvent]

  // OnRecordAuthWithOAuth2Request hook is triggered on each Record
  // OAuth2 sign-in/sign-up API request (after token exchange and before external provider linking).
  //
  // If [RecordAuthWithOAuth2RequestEvent.Record] is not set, then the OAuth2
  // request will try to create a new auth Record.
  //
  // To assign or link a different existing record model you can
  // change the [RecordAuthWithOAuth2RequestEvent.Record] field.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*core.RecordAuthWithOAuth2RequestEvent]

  // OnRecordAuthRefreshRequest hook is triggered on each Record
  // auth refresh API request (right before generating a new auth token).
  //
  // Could be used to additionally validate the request data or implement
  // completely different auth refresh behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAuthRefreshRequest(tags ...string) *hook.TaggedHook[*core.RecordAuthRefreshRequestEvent]

  // OnRecordRequestPasswordResetRequest hook is triggered on
  // each Record request password reset API request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different password reset behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestPasswordResetRequestEvent]

  // OnRecordConfirmPasswordResetRequest hook is triggered on
  // each Record confirm password reset API request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*core.RecordConfirmPasswordResetRequestEvent]

  // OnRecordRequestVerificationRequest hook is triggered on
  // each Record request verification API request.
  //
  // Could be used to additionally validate the loaded request data or implement
  // completely different verification behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordRequestVerificationRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestVerificationRequestEvent]

  // OnRecordConfirmVerificationRequest hook is triggered on each
  // Record confirm verification API request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*core.RecordConfirmVerificationRequestEvent]

  // OnRecordRequestEmailChangeRequest hook is triggered on each
  // Record request email change API request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different request email change behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestEmailChangeRequestEvent]

  // OnRecordConfirmEmailChangeRequest hook is triggered on each
  // Record confirm email change API request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*core.RecordConfirmEmailChangeRequestEvent]

  // OnRecordRequestOTPRequest hook is triggered on each Record
  // request OTP API request.
  //
  // [RecordCreateOTPRequestEvent.Record] could be nil if no matching identity is found, allowing
  // you to manually create or locate a different Record model (by reassigning [RecordCreateOTPRequestEvent.Record]).
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordRequestOTPRequest(tags ...string) *hook.TaggedHook[*core.RecordCreateOTPRequestEvent]

  // OnRecordAuthWithOTPRequest hook is triggered on each Record
  // auth with OTP API request.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordAuthWithOTPRequest(tags ...string) *hook.TaggedHook[*core.RecordAuthWithOTPRequestEvent]

  // ---------------------------------------------------------------
  // Record CRUD API event hooks
  // ---------------------------------------------------------------

  // OnRecordsListRequest hook is triggered on each API Records list request.
  //
  // Could be used to validate or modify the response before returning it to the client.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordsListRequest(tags ...string) *hook.TaggedHook[*core.RecordsListRequestEvent]

  // OnRecordViewRequest hook is triggered on each API Record view request.
  //
  // Could be used to validate or modify the response before returning it to the client.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordViewRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestEvent]

  // OnRecordCreateRequest hook is triggered on each API Record create request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordCreateRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestEvent]

  // OnRecordUpdateRequest hook is triggered on each API Record update request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordUpdateRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestEvent]

  // OnRecordDeleteRequest hook is triggered on each API Record delete request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different delete behavior.
  //
  // If the optional "tags" list (Collection ids or names) is specified,
  // then all event handlers registered via the created hook will be
  // triggered and called only if their event data origin matches the tags.
  OnRecordDeleteRequest(tags ...string) *hook.TaggedHook[*core.RecordRequestEvent]

  // ---------------------------------------------------------------
  // Collection API event hooks
  // ---------------------------------------------------------------

  // OnCollectionsListRequest hook is triggered on each API Collections list request.
  //
  // Could be used to validate or modify the response before returning it to the client.
  OnCollectionsListRequest() *hook.Hook[*core.CollectionsListRequestEvent]

  // OnCollectionViewRequest hook is triggered on each API Collection view request.
  //
  // Could be used to validate or modify the response before returning it to the client.
  OnCollectionViewRequest() *hook.Hook[*core.CollectionRequestEvent]

  // OnCollectionCreateRequest hook is triggered on each API Collection create request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  OnCollectionCreateRequest() *hook.Hook[*core.CollectionRequestEvent]

  // OnCollectionUpdateRequest hook is triggered on each API Collection update request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different persistence behavior.
  OnCollectionUpdateRequest() *hook.Hook[*core.CollectionRequestEvent]

  // OnCollectionDeleteRequest hook is triggered on each API Collection delete request.
  //
  // Could be used to additionally validate the request data or implement
  // completely different delete behavior.
  OnCollectionDeleteRequest() *hook.Hook[*core.CollectionRequestEvent]

  // OnCollectionsImportRequest hook is triggered on each API
  // collections import request.
  //
  // Could be used to additionally validate the imported collections or
  // to implement completely different import behavior.
  OnCollectionsImportRequest() *hook.Hook[*core.CollectionsImportRequestEvent]

  // ---------------------------------------------------------------
  // Batch API event hooks
  // ---------------------------------------------------------------

  // OnBatchRequest hook is triggered on each API batch request.
  //
  // Could be used to additionally validate or modify the submitted batch requests.
  OnBatchRequest() *hook.Hook[*core.BatchRequestEvent]
}
