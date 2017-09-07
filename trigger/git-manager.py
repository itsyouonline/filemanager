#!/usr/bin/env python3
import sys
import os
import subprocess
from pathlib import Path

class FilemanagerTrigger:
    def __init__(self):
        self.savedpath = None

        if not os.environ.get("TRIGGER"):
            print("[-] no trigger set, nothing to do")
            return None

        self.action = os.environ["TRIGGER"]
        self.username = os.environ["USERNAME"]
        self.realname = os.environ["REALNAME"]
        self.email = os.environ["USEREMAIL"]
        self.directory = os.environ["ROOT"]

        self.filetarget = ""
        self.filesource = os.environ["FILE"]
        self.fullsource = os.path.join(self.directory, self.filesource[1:])

        if os.environ.get("DESTINATION"):
            self.filetarget = os.environ["DESTINATION"]
            self.fulltarget = os.path.join(self.directory, self.filetarget[1:])

        print("[+] action: %s" % self.action)
        print("[+] owner: %s (%s, %s)" % (self.username, self.realname, self.email))
        print("[+] target: %s (file: %s)" % (self.directory, self.filesource))

    def process(self):
        if self.action == "before_save":
            return self.before_save()

        if self.action == "after_save":
            return self.after_save()

        if self.action == "before_publish":
            return self.before_publish()

        if self.action == "after_publish":
            return self.after_publish()

        if self.action == "after_copy":
            return self.after_copy()

        if self.action == "after_rename":
            return self.after_rename()

        if self.action == "after_upload":
            return self.after_upload()

        if self.action == "after_delete":
            return self.after_delete()

        print("[-] unknown trigger: %s" % self.action)
        return False

    def move(self, path):
        self.savedpath = os.getcwd()
        os.chdir(path)

    def restore(self):
        os.chdir(self.savedpath)
        self.savedpath = None

    def repository(self, fullpath):
        return os.path.dirname(fullpath)

    """
    Triggers (actions) implementation
    """
    # this trigger is fired before any changes
    def before_save(self):
        pass

    # this trigger is fired after file change (new content)
    def after_save(self):
        print("[+] updated: %s" % self.fullsource)
        repository = self.repository(self.fullsource)
        targetfile = os.path.basename(self.fullsource)

        author = "%s <%s>" % (self.realname, self.email)
        message = "Update %s [by %s]" % (targetfile, self.username)

        self.move(repository)
        subprocess.run(["git", "add", targetfile])
        subprocess.run(["git", "commit", "--author", author, "-m", message])
        subprocess.run(["git", "push", "origin", "master"])
        self.restore()

    def before_publish(self):
        pass

    def after_publish(self):
        pass

    def after_copy(self):
        print("[+] copied: %s -> %s" % (self.fullsource, self.fulltarget))
        repository = self.repository(self.fulltarget)
        targetfile = os.path.basename(self.fulltarget)

        author = "%s <%s>" % (self.realname, self.email)
        message = "Copy %s [by %s]" % (targetfile, self.username)

        self.move(repository)
        subprocess.run(["git", "add", targetfile])
        subprocess.run(["git", "commit", "--author", author, "-m", message])
        subprocess.run(["git", "push", "origin", "master"])
        self.restore()

    def after_rename(self):
        print("[+] renamed: %s -> %s" % (self.fullsource, self.fulltarget))
        # need to support cross-repository

    def after_upload(self):
        print("[+] uploaded: %s" % self.fullsource)
        repository = self.repository(self.fullsource)
        targetfile = os.path.basename(self.fullsource)

        author = "%s <%s>" % (self.realname, self.email)
        message = "Create %s [by %s]" % (targetfile, self.username)

        self.move(repository)
        subprocess.run(["git", "add", targetfile])
        subprocess.run(["git", "commit", "--author", author, "-m", message])
        subprocess.run(["git", "push", "origin", "master"])
        self.restore()

    def after_delete(self):
        print("[+] deleted: %s" % self.fullsource)
        repository = self.repository(self.fullsource)
        targetfile = os.path.basename(self.fullsource)

        if self.fullsource.endswith("/"):
            repository = str(Path(repository).parent)

        author = "%s <%s>" % (self.realname, self.email)
        message = "Delete %s [by %s]" % (targetfile, self.username)

        self.move(repository)
        subprocess.run(["git", "add", "-u", targetfile])
        subprocess.run(["git", "commit", "--author", author, "-m", message])
        subprocess.run(["git", "push", "origin", "master"])
        self.restore()

if __name__ == '__main__':
    fm = FilemanagerTrigger()
    fm.process()
