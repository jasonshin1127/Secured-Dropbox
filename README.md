# Secured-Dropbox
Dropbox-like file-sharing system that incorporates cryptographic primitives

Users are allowed to take the following actions:
- Authenticate with a username and password;
- Save files to the server;
- Load saved files from the server;
- Overwrite saved files on the server;
- Append to saved files on the server;
- Share saved files with other users; and
- Revoke access to previously shared files.

Prevent attacks like:

1. Datastore Adversary:
The Datastore is an untrusted service hosted on a server and network controlled by an adversary. The adversary can view and record the content and metadata of all requests (set/get/delete) to the Datastore API. This allows the adversary to know who stored which key-value entry, when, and what the contents are.
2. Revoked User Adversary:
Malicious users may try to perform operations on arbitrary files by utilizing the request/response information that they recorded before their access was revoked. All writes to Datastore made by a user in an attempt to modify file content or re-acquire access to file are malicious actions.

Design Requirements:

1. Usernames / Passwords
(Usernames)
The client SHOULD assume that each user has a unique username.
Usernames are case-sensitive: Bob and bob are different users.
The client SHOULD support usernames of any length greater than zero.
(Passwords)
The client MUST NOT assume each user has a unique password. Like the real world, users may happen to choose the same password.
The client MAY assume each user’s password generally is a good source of entropy. However, the attackers possess a precomputed lookup table containing hashes of common passwords downloaded from the internet.
The client SHOULD support passwords length greater than or equal to zero.

2. User Sessions
The client application MUST allow many different users to use the application at the same time. For example, Bob and Alice can each run the client application on their own devices at the same time.

The client MUST support a single user having multiple active sessions at the same time. All file changes MUST be reflected in all current user sessions immediately (i.e. without terminating the current session and re-authenticating).

For example:

Bob runs the client application on his laptop and calls InitUser() to create session bobLaptop.
Bob wants to run the client application on his tablet, so he calls GetUser on his tablet to get bobTablet.
Using bobLaptop, Bob stores a file file1.txt. Session bobTablet must be able to download file1.txt.
Using bobTablet, Bob appends to file1.txt. Session bobLaptop must be able to download the updated version.
Using bobLaptop, Bob accepts an invitation to access a file and calls the file file2.txt in his personal namespace. Bob must be able to load the corresponding file2.txt using bobTablet.
The client DOES NOT need to support concurrency. Globally, across all users and user-sessions, all operations in the client application will be done serially.


3. Cryptography and Keys
Each public key SHOULD be used for a single purpose, which means that each user is likely to have more than one public key.
A single user MAY have multiple entries in Keystore. However, the number of keys in Keystore per user MUST be a small constant; it MUST NOT depend on the number of files stored or length of any file, how many users a file has been shared with, or the number of users already in the system.
The following SHOULD be avoided because they are dangerous design patterns that often lead to subtle vulnerabilities:
- Reusing the same key for multiple purposes.
- Authenticate-then-encrypt.
- Decrypt-then-verify.


4. No Persistent Local State
The client MUST NOT save any data to the local file system. If the client is restarted, it must be able to pick up where it left off given only a username and password. Any data requiring persistent storage MUST be stored in either Keystore or Datastore.


5. Files
Any breach of IND-CPA security constitutes loss of confidentiality.
The client MUST ensure confidentiality and integrity of file contents and file sharing invitations.
The client MUST ensure the integrity of filenames.
The client MUST prevent adversaries from learning filenames and the length of filenames. The client MAY use and store filenames in a deterministic manner.
The client MUST prevent the revoked user adversary from learning anything about future writes or appends to the file after their access has been revoked.
The client MAY leak any information except filenames, lengths of filenames, file contents, and file sharing invitations. For example, the client design MAY leak the size of file contents or the number of files associated with a user.
Filenames MAY be any length, including zero (empty string).
The client MUST NOT assume that filenames are globally unique. For example, user bob can have a file named foo.txt and user alice can have a file named foo.txt. The client MUST keep each user’s file namespace independent from one another.

6. Sharing and Revocation
<img width="470" alt="image" src="https://github.com/jasonshin1127/Secured-Dropbox/assets/101506840/dee1c076-1c0e-4a2b-aaf8-ed466d6b3c06">

The client MUST enforce authorization for all files. The only users who are authorized to access a file using the client include: (a) the owner of the file; and (b) users who have accepted an invitation to access the file and that access has not been revoked.

The client MUST allow any user who is authorized to access the file to take the following actions on the file:

- Read file contents (LoadFile()).
- Overwrite file contents (StoreFile()).
- Append additional contents to the file (AppendToFile()).
- Share the file with other users (CreateInvitation()).
For example, all of the users listed in Figure 1 are authorized to take the listed actions on the file.

Changes to the contents of a file MUST be accessible by all users who are authorized to access the file.

The client MUST enforce that there is only a single copy of a file. Sharing the file MAY NOT create a copy of the file.

The client MUST ensure the confidentiality and integrity of the secure file share invitations created by CreateInvitation().

The client MAY assume that CreateInvitation() will never be called on recipients who are already authorized to access the file.

The client MAY assume that CreateInvitation() will never be called in a situation where the sharing action will result in an ill-formed tree structure. In a well-formed tree structure, each non-root node has a single parent and there are no cycles. For example, in Figure 1, Nilufar’s attempts to share with Marco or Alice in steps 5 and 6 would be undefined behavior, since both users are already authorized to access the file, and the sharing actions would create an ill-formed tree structure.

The client MUST enforce that the file owner is able to revoke access from users who they directly shared the file with.

Any other revocation (i.e. owners revoking users who they did not directly share with, or revocations by non-owners) is undefined behavior and will not be tested.

For example, in Figure 1, Alice is the only user who MUST be able to revoke access, and she MUST be able to revoke access from Bob and Nilufar. If any user other than Alice attempts to revoke access, or Alice attempts to revoke access from any user other than Bob orNilufar, this is undefined behavior and will not be tested.

When the owner revokes a user’s access, the client MUST enforce that any other users with whom the revoked user previously shared the file also lose access.

For example, in Figure 1, if Alice revokes access from Bob, then all of the following users MUST lose access: Bob, Olga, and Marco. As the file owner, Alice always maintains access. Nilufar maintains access because Bob did not grant Nilufar access to the file (Alice did).

The client MUST prevent any revoked user from using the client API to take any action on the file from which their access was revoked. However, recall from Threat Model that a revoked user may become malicious and use the Datastore API directly.

Re-sharing a file with a revoked user is undefined behavior and will not be tested.



